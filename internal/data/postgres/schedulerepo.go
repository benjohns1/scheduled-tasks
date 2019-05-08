package postgres

import (
	"database/sql"
	"fmt"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
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
	for rows.Next() {
		sd, err := parseScheduleRow(rows)
		if err != nil {
			return nil, usecase.NewError(usecase.ErrUnknown, "error parsing schedule row: %v", err)
		}

		// Add to cache
		r.schedules[sd.ScheduleID] = sd.Schedule
	}

	return r.schedules, nil
}

func scheduleSelectClause() (selectClause string) {
	return "SELECT id, paused FROM schedule"
}

func parseScheduleRow(r scannable) (sd usecase.ScheduleData, err error) {

	sd = usecase.ScheduleData{}

	// Scan into row data structure
	var row struct {
		id     int64
		paused bool
	}
	err = r.Scan(&row.id, &row.paused)
	if err != nil {
		return
	}

	// Construct schedule entity
	sd.Schedule, err = schedule.NewRaw(row.paused)
	sd.ScheduleID = usecase.ScheduleID(row.id)

	return
}

// Add adds a schedule to the persisence layer
func (r *ScheduleRepo) Add(s *schedule.Schedule) (usecase.ScheduleID, usecase.Error) {
	q := "INSERT INTO schedule (paused) VALUES ($1) RETURNING id"
	var id usecase.ScheduleID
	err := r.db.QueryRow(q, s.Paused()).Scan(&id)
	if err != nil {
		return 0, usecase.NewError(usecase.ErrUnknown, "error inserting new schedule: %v", err)
	}

	r.schedules[id] = s

	return id, nil
}

// Update updates a schedule's persistent data to the given aggregate values
func (r *ScheduleRepo) Update(id usecase.ScheduleID, s *schedule.Schedule) usecase.Error {
	q := "UPDATE schedule SET paused = $1 WHERE id = $2 RETURNING id"
	rows, err := r.db.Query(q, s.Paused(), id)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error updating schedule id %d: %v", id, err)
	}
	if !rows.Next() {
		return usecase.NewError(usecase.ErrRecordNotFound, "no schedule found for id = %v", id)
	}

	r.schedules[id] = s

	return nil
}
