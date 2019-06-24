package usecase

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
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
	GetAllScheduled() (map[ScheduleID]*schedule.Schedule, Error)
	Add(*schedule.Schedule) (ScheduleID, Error)
	Update(ScheduleID, *schedule.Schedule) Error
}

// GetSchedule returns a single schedule
func GetSchedule(r ScheduleRepo, id ScheduleID) (*ScheduleData, Error) {
	s, err := r.Get(id)
	if err != nil {
		return nil, err.Prefix("error getting schedule id %v", id)
	}

	if !s.IsValid() {
		return nil, NewError(ErrRecordNotFound, "schedule id %d not found", id)
	}

	return &ScheduleData{ScheduleID: id, Schedule: s}, nil
}

// ListSchedules returns all schedules
func ListSchedules(r ScheduleRepo) (map[ScheduleID]*schedule.Schedule, Error) {
	ss, err := r.GetAll()
	if err != nil {
		return nil, err.Prefix("error listing all schedules")
	}

	list := make(map[ScheduleID]*schedule.Schedule)
	for id, s := range ss {
		if !s.IsValid() {
			continue
		}
		list[id] = s
	}

	return list, nil
}

// AddSchedule adds a new schedule
func AddSchedule(r ScheduleRepo, s *schedule.Schedule, checkSchedule chan<- bool) (ScheduleID, Error) {
	id, err := r.Add(s)
	if err != nil {
		return id, err.Prefix("error adding schedule")
	}
	select {
	case checkSchedule <- true:
	default:
	}
	return id, nil
}

// PauseSchedule pauses the schedule
func PauseSchedule(r ScheduleRepo, id ScheduleID, checkSchedule chan<- bool) Error {

	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %d to pause", id)
	}

	if !s.IsValid() {
		return NewError(ErrRecordNotFound, "schedule id %d not found", id)
	}

	s.Pause()

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to pause", id)
	}
	select {
	case checkSchedule <- true:
	default:
	}
	return nil
}

// UnpauseSchedule unpauses the schedule
func UnpauseSchedule(r ScheduleRepo, id ScheduleID, checkSchedule chan<- bool) Error {

	s, err := r.Get(id)
	if err != nil {
		return err.Prefix("error retrieving schedule id %d to unpause", id)
	}

	if !s.IsValid() {
		return NewError(ErrRecordNotFound, "schedule id %d not found", id)
	}

	s.Unpause()

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to unpause", id)
	}
	select {
	case checkSchedule <- true:
	default:
	}
	return nil
}

// RemoveSchedule removes a schedule
func RemoveSchedule(r ScheduleRepo, id ScheduleID, checkSchedule chan<- bool) Error {
	s, ucErr := r.Get(id)
	if ucErr != nil {
		return ucErr.Prefix("error retrieving schedule id %d to remove", id)
	}

	err := s.Remove()
	if err != nil {
		return NewError(ErrUnknown, "error removing schedule id %d", id)
	}

	ucErr = r.Update(id, s)
	if err != nil {
		return ucErr.Prefix("error attempting to remove schedule id %d", id)
	}
	select {
	case checkSchedule <- true:
	default:
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

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to add recurring task", id)
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

	err = r.Update(id, s)
	if err != nil {
		return err.Prefix("error updating schedule id %d attempting to remove recurring task", id)
	}

	return nil
}
