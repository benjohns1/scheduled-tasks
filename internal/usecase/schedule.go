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
	Add(*schedule.Schedule) (TaskID, Error)
	Update(ScheduleID, *schedule.Schedule) Error
}

// GetSchedule returns a single schedule
func GetSchedule(r ScheduleRepo, id ScheduleID) (*ScheduleData, Error) {
	return nil, NewError(ErrUnknown, "not implemented")
}

// ListSchedules returns all schedules
func ListSchedules(r ScheduleRepo) (map[ScheduleID]*schedule.Schedule, Error) {
	return map[ScheduleID]*schedule.Schedule{}, NewError(ErrUnknown, "not implemented")
}

// AddSchedule adds a new schedule
func AddSchedule(r ScheduleRepo, s *schedule.Schedule) (ScheduleID, Error) {
	return 0, NewError(ErrUnknown, "not implemented")
}

// PauseSchedule pauses the schedule
func PauseSchedule(r ScheduleRepo, id ScheduleID) Error {
	return NewError(ErrUnknown, "not implemented")
}

// UnpauseSchedule unpauses the schedule
func UnpauseSchedule(r ScheduleRepo, id ScheduleID) Error {
	return NewError(ErrUnknown, "not implemented")
}

// AddRecurringTask adds a new recurring task to the schedule
func AddRecurringTask(r ScheduleRepo, id ScheduleID, rt *schedule.RecurringTask) Error {
	return NewError(ErrUnknown, "not implemented")
}

// RemoveRecurringTask removes the recurring task at the specified index from the schedule
func RemoveRecurringTask(r ScheduleRepo, id ScheduleID, taskIndex int) Error {
	return NewError(ErrUnknown, "not implemented")
}
