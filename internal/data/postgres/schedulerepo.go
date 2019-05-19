package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"

	"github.com/lib/pq"
)

// ScheduleRepo handles persisting task data and maintaining an in-memory cache
type ScheduleRepo struct {
	db        *sql.DB
	schedules map[usecase.ScheduleID]*schedule.Schedule
}

// NewScheduleRepo instantiates a new ScheduleRepo
func NewScheduleRepo(conn DBConn) (repo *ScheduleRepo, err error) {

	if conn.DB == nil {
		return nil, fmt.Errorf("DB connection is nil")
	}

	return &ScheduleRepo{db: conn.DB, schedules: make(map[usecase.ScheduleID]*schedule.Schedule)}, nil
}

// Get retrieves a schedule aggregate, given its persistent ID
func (r *ScheduleRepo) Get(id usecase.ScheduleID) (*schedule.Schedule, usecase.Error) {

	// Try to retrieve from cache
	s, ok := r.schedules[id]
	if ok {
		return s, nil
	}

	// Retrieve from DB
	query := fmt.Sprintf("%s WHERE id = $1", scheduleSelectClause())
	row := r.db.QueryRow(query, id)
	sd, err := parseScheduleRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.NewError(usecase.ErrRecordNotFound, "no task found with id = %v", id)
		}
		return nil, usecase.NewError(usecase.ErrUnknown, "error parsing schedule id %d: %v", id, err)
	}

	// Get recurring tasks from DB
	rts, err := r.getRecurringTasks([]usecase.ScheduleID{id})
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving recurring tasks for schedule id %v", id)
	}
	for _, rt := range rts[id] {
		sd.Schedule.AddTask(rt)
	}

	// Add to cache
	r.schedules[sd.ScheduleID] = sd.Schedule

	return sd.Schedule, nil
}

// GetAll retrieves all schedules
func (r *ScheduleRepo) GetAll() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {
	// Retrieve from DB
	rows, err := r.db.Query(scheduleSelectClause())
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving all schedules: %v", err)
	}

	scheds := map[usecase.ScheduleID]*schedule.Schedule{}
	sids := []usecase.ScheduleID{}
	for rows.Next() {
		sd, err := parseScheduleRow(rows)
		if err != nil {
			return nil, usecase.NewError(usecase.ErrUnknown, "error parsing schedule row: %v", err)
		}

		scheds[sd.ScheduleID] = sd.Schedule
		sids = append(sids, sd.ScheduleID)
	}

	// Get recurring tasks from DB
	allTasks, err := r.getRecurringTasks(sids)
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving recurring tasks for all schedules: %v", err)
	}
	for sid, rts := range allTasks {
		for _, rt := range rts {
			scheds[sid].AddTask(rt)
		}
	}

	// Replace cache
	r.schedules = scheds

	return r.schedules, nil
}

func scheduleSelectClause() (selectClause string) {
	return "SELECT id, paused, frequency_offset, frequency_interval, frequency_time_period, frequency_at_minutes FROM schedule"
}

func parseScheduleRow(r scannable) (sd usecase.ScheduleData, err error) {

	sd = usecase.ScheduleData{}

	// Scan into row data structure
	var row struct {
		id          int64
		fOffset     int
		fInterval   int
		fTimePeriod schedule.TimePeriod
		fAtMinutes  []sql.NullInt64
		paused      bool
	}
	err = r.Scan(&row.id, &row.paused, &row.fOffset, &row.fInterval, &row.fTimePeriod, pq.Array(&row.fAtMinutes))
	if err != nil {
		return
	}

	// Construct frequency value
	f, err := schedule.NewRawFrequency(row.fOffset, row.fInterval, row.fTimePeriod, toIntSlice(row.fAtMinutes), nil, nil, nil)
	if err != nil {
		return
	}

	// Construct schedule entity
	sd.Schedule = schedule.NewRaw(f, row.paused, []schedule.RecurringTask{})
	sd.ScheduleID = usecase.ScheduleID(row.id)

	return
}

func toIntSlice(sqlSlice []sql.NullInt64) []int {
	if sqlSlice == nil {
		return nil
	}
	intSlice := make([]int, len(sqlSlice))
	for i, item := range sqlSlice {
		if !item.Valid {
			continue
		}
		intSlice[i] = int(item.Int64)
	}
	return intSlice
}

// Add adds a schedule to the persisence layer
func (r *ScheduleRepo) Add(s *schedule.Schedule) (usecase.ScheduleID, usecase.Error) {
	q := "INSERT INTO schedule (paused, frequency_offset, frequency_interval, frequency_time_period, frequency_at_minutes) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var id usecase.ScheduleID
	f := s.Frequency()
	err := r.db.QueryRow(q, s.Paused(), f.Offset(), f.Interval(), f.TimePeriod(), pq.Array(f.AtMinutes())).Scan(&id)
	if err != nil {
		return 0, usecase.NewError(usecase.ErrUnknown, "error inserting new schedule: %v", err)
	}
	rts := s.Tasks()
	if len(rts) > 0 {
		err := r.insertTasks(id, rts)
		if err != nil {
			return 0, usecase.NewError(usecase.ErrUnknown, "error inserting recurring tasks to schedule: %v", err)
		}
	}

	r.schedules[id] = s

	return id, nil
}

func parseRecurringTaskRow(r scannable) (id int64, sid usecase.ScheduleID, rt schedule.RecurringTask, err error) {

	rt = schedule.RecurringTask{}

	// Scan into row data structure
	var row struct {
		name        string
		description string
	}
	err = r.Scan(&id, &sid, &row.name, &row.description)
	if err != nil {
		return
	}

	// Construct recurring task value object
	rt = schedule.NewRecurringTask(row.name, row.description)
	return
}

func (r *ScheduleRepo) getRecurringTasks(sids []usecase.ScheduleID) (map[usecase.ScheduleID]map[int64]schedule.RecurringTask, error) {
	ts := map[usecase.ScheduleID]map[int64]schedule.RecurringTask{}
	if len(sids) <= 0 {
		return ts, nil
	}

	sidsString := make([]string, len(sids))
	for i, sid := range sids {
		sidsString[i] = strconv.Itoa(int(sid))
	}
	q := fmt.Sprintf("SELECT id, schedule_id, name, description FROM recurring_task WHERE schedule_id IN (%s)", strings.Join(sidsString, ","))
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tasks: %v", err)
	}
	for rows.Next() {
		id, sid, t, err := parseRecurringTaskRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error parsing task row: %v", err)
		}
		if _, ok := ts[sid]; !ok {
			ts[sid] = map[int64]schedule.RecurringTask{}
		}
		ts[sid][id] = t
	}
	return ts, nil
}

func (r *ScheduleRepo) insertTasks(sid usecase.ScheduleID, rts []schedule.RecurringTask) error {
	q := "INSERT INTO recurring_task (schedule_id, name, description) VALUES ($1, $2, $3) RETURNING id"
	var rtid int64
	for _, rt := range rts {
		err := r.db.QueryRow(q, sid, rt.Name(), rt.Description()).Scan(&rtid)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update updates a schedule's persistent data to the given aggregate values
func (r *ScheduleRepo) Update(id usecase.ScheduleID, s *schedule.Schedule) usecase.Error {
	q := "UPDATE schedule SET paused = $1, frequency_offset = $2, frequency_interval = $3, frequency_time_period = $4, frequency_at_minutes = $5 WHERE id = $6 RETURNING id"
	f := s.Frequency()
	rows, err := r.db.Query(q, s.Paused(), f.Offset(), f.Interval(), f.TimePeriod(), pq.Array(f.AtMinutes()), id)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error updating schedule id %d: %v", id, err)
	}
	if !rows.Next() {
		return usecase.NewError(usecase.ErrRecordNotFound, "no schedule found for id = %v", id)
	}

	r.schedules[id] = s

	return nil
}
