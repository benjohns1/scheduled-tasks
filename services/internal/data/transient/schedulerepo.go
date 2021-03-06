package transient

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
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
	s, ok := r.schedules[id]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no schedule with ID: %v", id)
	}
	return s, nil
}

// GetForUser retrieves a schedule entity for a user, given its persistent ID
func (r *ScheduleRepo) GetForUser(id usecase.ScheduleID, uid user.ID) (*schedule.Schedule, usecase.Error) {
	s, ok := r.schedules[id]
	if !ok || !uid.Equals(s.CreatedBy()) {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no schedule with ID: %v", id)
	}
	return s, nil
}

// GetAllScheduled retrieves all valid, unpaused schedules
func (r *ScheduleRepo) GetAllScheduled() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {
	scheds := map[usecase.ScheduleID]*schedule.Schedule{}
	for id, s := range r.schedules {
		if s.IsValid() && !s.Paused() {
			scheds[id] = s
		}
	}
	return scheds, nil
}

// GetAll retrieves all schedules
func (r *ScheduleRepo) GetAll() (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {

	return r.schedules, nil
}

// GetAllForUser retrieves all schedules created by the given user
func (r *ScheduleRepo) GetAllForUser(uid user.ID) (map[usecase.ScheduleID]*schedule.Schedule, usecase.Error) {
	ss := map[usecase.ScheduleID]*schedule.Schedule{}
	for id, s := range r.schedules {
		if s.IsValid() && uid.Equals(s.CreatedBy()) {
			ss[id] = s
		}
	}
	return ss, nil
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
