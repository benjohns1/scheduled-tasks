package scheduler

import (
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Offset is the offset added to the next run time
const Offset = 3 * time.Second

// Run starts the scheduler process
func Run(l Logger, c usecase.Clock, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo) (close chan<- bool, closed <-chan bool, next <-chan time.Time) {
	l.Printf("scheduler process starting")

	closeSignal := make(chan bool)
	onClosed := make(chan bool)
	nextRunTimes := make(chan time.Time)

	go func() {
		defer func() {
			onClosed <- true
		}()
		for {
			l.Printf("checking schedules")
			nextRecurrence, err := usecase.CheckSchedules(c, taskRepo, scheduleRepo)
			if err != nil {
				l.Printf("halting scheduler: %v", err)
				break
			}
			if nextRecurrence.IsZero() {
				l.Printf("no upcoming schedules, halting")
				break
			}

			// Sleep until next scheduled time + offset
			l.Printf("next run scheduled for %v + %v offset", nextRecurrence, Offset)
			nextRunTime := nextRecurrence.Add(Offset)
			wait := c.Until(nextRunTime)
			nextRunTimes <- nextRunTime

			// Listen for exit signal or until next schedule process
			select {
			case <-closeSignal:
				l.Printf("scheduler exiting")
				return
			case <-c.After(wait):
			}
		}
		l.Printf("scheduler process complete")
	}()

	return closeSignal, onClosed, nextRunTimes
}
