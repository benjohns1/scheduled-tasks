package core

import (
	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
)

type Schedule struct {
	frequency *schedule.Frequency
	Paused    bool
	tasks     []RecurringTask
}

type RecurringTask struct {
}

func NewSchedule() *Schedule {
	return &Schedule{frequency: schedule.NewFrequency()}
}
