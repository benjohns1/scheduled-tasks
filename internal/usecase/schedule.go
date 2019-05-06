package usecase

import (
	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
)

// ScheduleID is the persistent ID of the task
type ScheduleID int64

// ScheduleData contains application-level task info
type ScheduleData struct {
	ScheduleID ScheduleID
	Schedule   *schedule.Schedule
}

// ScheduleRepo defines the task repository interface required by use cases
type ScheduleRepo interface {
	Get(ScheduleID) (*schedule.Schedule, Error)
	GetAll() (map[ScheduleID]*schedule.Schedule, Error)
	Add(*schedule.Schedule) (ScheduleID, Error)
	Update(ScheduleID, *schedule.Schedule) Error
}

// GetSchedule returns a single schedule
func GetSchedule(r ScheduleRepo, id ScheduleID) (*ScheduleData, Error) {
	s, err := r.Get(id)
	if err != nil {
		return nil, err.Prefix("error getting schedule id %v", id)
	}
	return &ScheduleData{ScheduleID: id, Schedule: s}, nil
}

// ListSchedules returns all schedules
func ListSchedules(r ScheduleRepo) (map[ScheduleID]*schedule.Schedule, Error) {
	ss, err := r.GetAll()
	if err != nil {
		return nil, err.Prefix("error listing all schedules")
	}
	return ss, nil
}

// AddSchedule adds a new schedule
func AddSchedule(r ScheduleRepo, s *schedule.Schedule) (ScheduleID, Error) {
	id, err := r.Add(s)
	if err != nil {
		return id, err.Prefix("error adding schedule")
	}
	return id, nil
}

// PauseSchedule pauses the schedule
func PauseSchedule(r ScheduleRepo, id ScheduleID) Error {

	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %d to pause", id)
	}

	s.Pause()

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to pause", id)
	}
	return nil
}

// UnpauseSchedule unpauses the schedule
func UnpauseSchedule(r ScheduleRepo, id ScheduleID) Error {

	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %d to unpause", id)
	}

	s.Unpause()

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to unpause", id)
	}
	return nil
}

// AddRecurringTask adds a new recurring task to the schedule
func AddRecurringTask(r ScheduleRepo, id ScheduleID, rt schedule.RecurringTask) Error {

	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %v to add recurring task", id)
	}

	if e := s.AddTask(rt); e != nil {
		return NewError(ErrDuplicateRecord, "can't add recurring task: duplicate found for schedule id %v", id)
	}

	return nil
}

// RemoveRecurringTask removes the recurring task at the specified index from the schedule
func RemoveRecurringTask(r ScheduleRepo, id ScheduleID, rt schedule.RecurringTask) Error {
	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %v to remove recurring task", id)
	}

	if e := s.RemoveTask(rt); e != nil {
		return NewError(ErrRecordNotFound, "can't remove recurring task from schedule id %v", id)
	}

	return nil
}
