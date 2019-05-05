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
