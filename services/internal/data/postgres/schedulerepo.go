package postgres

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"

	"github.com/lib/pq"
)

// ScheduleRepo persists schedule data in a PostgreSQL DB
type ScheduleRepo struct {
	db *sql.DB
}

// NewScheduleRepo instantiates a new ScheduleRepo
func NewScheduleRepo(conn DBConn) (repo *ScheduleRepo, err error) {

	if conn.DB == nil {
		return nil, fmt.Errorf("DB connection is nil")
	}

	return &ScheduleRepo{db: conn.DB}, nil
}

// Get retrieves a schedule aggregate, given its persistent ID
func (r *ScheduleRepo) Get(id usecase.ScheduleID) (*schedule.Schedule, usecase.Error) {

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

	return sd.Schedule, nil
}

// GetAll retrieves all schedules
func (r *ScheduleRepo) GetAll() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {
	return r.getAllWhere("")
}

// GetAllScheduled retrieves all unpaused schedules that haven't been removed
func (r *ScheduleRepo) GetAllScheduled() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {
	return r.getAllWhere("paused = FALSE AND removed_time = $1", time.Time{})
}

func (r *ScheduleRepo) getAllWhere(whereClause string, params ...interface{}) (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {

	q := scheduleSelectClause()
	if whereClause != "" {
		q = fmt.Sprintf("%v WHERE %v", q, whereClause)
	}

	// Retrieve from DB
	rows, err := r.db.Query(q, params...)
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving all schedules: %v", err)
	}
	defer rows.Close()

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

	return scheds, nil
}

func scheduleSelectClause() (selectClause string) {
	return "SELECT id, paused, last_checked, removed_time, frequency_offset, frequency_interval, frequency_time_period, frequency_at_minutes FROM schedule"
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
		lastChecked *string
		removed     *string
	}
	err = r.Scan(&row.id, &row.paused, &row.lastChecked, &row.removed, &row.fOffset, &row.fInterval, &row.fTimePeriod, pq.Array(&row.fAtMinutes))
	if err != nil {
		return
	}

	// Construct frequency value
	f, err := schedule.NewRawFrequency(row.fOffset, row.fInterval, row.fTimePeriod, toIntSlice(row.fAtMinutes), nil, nil, nil)
	if err != nil {
		return
	}

	lastChecked, err := time.Parse(dbTimeFormat, *row.lastChecked)
	if err != nil {
		lastChecked = time.Time{}
	}
	removed, err := time.Parse(dbTimeFormat, *row.removed)
	if err != nil {
		removed = time.Time{}
	}

	// Construct schedule entity
	sd.Schedule = schedule.NewRaw(f, row.paused, lastChecked, []schedule.RecurringTask{}, removed)
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
	q := "INSERT INTO schedule (paused, last_checked, removed_time, frequency_offset, frequency_interval, frequency_time_period, frequency_at_minutes) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	var id usecase.ScheduleID
	f := s.Frequency()
	err := r.db.QueryRow(q, s.Paused(), s.LastChecked(), s.RemovedTime(), f.Offset(), f.Interval(), f.TimePeriod(), pq.Array(f.AtMinutes())).Scan(&id)
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
	defer rows.Close()
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

func (r *ScheduleRepo) clearTasks(sid usecase.ScheduleID) error {
	q := "DELETE FROM recurring_task WHERE schedule_id = $1"
	_, err := r.db.Exec(q, sid)
	if err != nil {
		return fmt.Errorf("error clearing all tasks from recurring_task table: %v", err)
	}
	return nil
}

// Update updates a schedule's persistent data to the given aggregate values
func (r *ScheduleRepo) Update(id usecase.ScheduleID, s *schedule.Schedule) usecase.Error {

	// Update schedule row
	q := "UPDATE schedule SET paused = $1, last_checked = $2, removed_time = $3, frequency_offset = $4, frequency_interval = $5, frequency_time_period = $6, frequency_at_minutes = $7 WHERE id = $8 RETURNING id"
	f := s.Frequency()
	rows, err := r.db.Query(q, s.Paused(), s.LastChecked(), s.RemovedTime(), f.Offset(), f.Interval(), f.TimePeriod(), pq.Array(f.AtMinutes()), id)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error updating schedule id %d: %v", id, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return usecase.NewError(usecase.ErrRecordNotFound, "no schedule found for id = %v", id)
	}

	// Check if any tasks need to be modified
	rts, err := r.getRecurringTasks([]usecase.ScheduleID{id})
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error retrieving recurring tasks for schedule id %v: %v", id, err)
	}
	newRts := s.Tasks()
	if AnyTasksModified(rts[id], newRts) {
		err := r.replaceTasks(id, newRts)
		if err != nil {
			return usecase.NewError(usecase.ErrUnknown, "error updating recurring tasks for schedule id %v: %v", id, err)
		}
	}

	return nil
}

// AnyTasksModified returns whether the map of recurring tasks contains all entries in the slice of recurring tasks
func AnyTasksModified(as map[int64]schedule.RecurringTask, bs []schedule.RecurringTask) bool {
	if len(as) != len(bs) {
		return true
	}
	usedIndices := make(map[int]bool)
	for _, at := range as {
		match := false
		for i, bt := range bs {
			if _, used := usedIndices[i]; used {
				continue
			}
			if at.Equal(bt) {
				match = true
				usedIndices[i] = true
				break
			}
		}
		if !match {
			return true
		}
	}
	return false
}

func (r *ScheduleRepo) replaceTasks(id usecase.ScheduleID, rts []schedule.RecurringTask) error {
	// Modify tasks by clearing and reinserting all
	// @TODO: determine which specific tasks need updating and only update those
	err := r.clearTasks(id)
	if err != nil {
		return fmt.Errorf("error clearing recurring tasks: %v", err)
	}
	if len(rts) > 0 {
		err := r.insertTasks(id, rts)
		if err != nil {
			return fmt.Errorf("error inserting recurring tasks: %v", err)
		}
	}
	return nil
}
