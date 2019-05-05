package transient

import (
	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// ScheduleRepo maintains an in-memory cache of tasks
type ScheduleRepo struct {
	lastID    int
	schedules map[usecase.ScheduleID]*schedule.Schedule
}

// NewScheduleRepo instantiates a new TaskRepo
func NewScheduleRepo() *ScheduleRepo {
	return &ScheduleRepo{schedules: make(map[usecase.ScheduleID]*schedule.Schedule)}
}

// Get retrieves a schedule entity, given its persistent ID
func (r *ScheduleRepo) Get(id usecase.ScheduleID) (*schedule.Schedule, usecase.Error) {

	// Try to retrieve from cache
	s, ok := r.schedules[id]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no schedule with ID: %v", id)
	}
	return s, nil
}

// GetAll retrieves all schedules
func (r *ScheduleRepo) GetAll() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {

	return r.schedules, nil
}

// Add adds a task to the persisence layer
func (r *ScheduleRepo) Add(s *schedule.Schedule) (usecase.ScheduleID, usecase.Error) {
	r.lastID++
	id := usecase.ScheduleID(r.lastID)
	r.schedules[id] = s

	return id, nil
}

// Update updates a task's persistent data to the given entity values
func (r *ScheduleRepo) Update(id usecase.ScheduleID, s *schedule.Schedule) usecase.Error {

	_, ok := r.schedules[id]
	if !ok {
		return usecase.NewError(usecase.ErrRecordNotFound, "no schedule with ID %v", id)
	}

	r.schedules[id] = s

	return nil
}
