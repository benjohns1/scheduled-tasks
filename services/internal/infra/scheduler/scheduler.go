package scheduler

import (
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Offset is the offset added to the next run time
const Offset = 3 * time.Second

// DefaultWait is the default amount of time for the process to wait before automatically checking the schedule, if there's no upcoming recurrences
const DefaultWait = 7 * 24 * time.Hour

// Run starts the scheduler process
func Run(l Logger, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo, nextRun chan time.Time) (close chan<- bool, check chan<- bool, closed <-chan bool) {
	l.Printf("scheduler process starting")

	checkSignal := make(chan bool)
	closeSignal := make(chan bool)
	onClosed := make(chan bool)

	go func() {
		defer func() {
			select {
			case onClosed <- true:
			default:
			}
		}()
		for {
			l.Printf("checking schedules")
			nextRecurrence, err := usecase.CheckSchedules(taskRepo, scheduleRepo)
			if err != nil {
				l.Printf("error checking schedules: %v", err)
			}
			if nextRecurrence.IsZero() {
				l.Printf("no upcoming schedules, setting default wait to check schedule in %v from now", DefaultWait)
				nextRecurrence = clock.Now().Add(DefaultWait)
			}

			// Sleep until next scheduled time + offset
			l.Printf("next run scheduled for %v + %v offset", nextRecurrence, Offset)
			nextRunTime := nextRecurrence.Add(Offset)

			// Notify receivers of next runtime
			if nextRun != nil {
				nextRun <- nextRunTime
			}

			wait := clock.Until(nextRunTime)
			if wait <= 0 {
				wait = 1
			}

			// Listen for exit signal, check signal, or until next recurrence is ready
			select {
			case <-closeSignal:
				l.Printf("scheduler exiting")
				return
			case <-checkSignal:
			case <-clock.After(wait):
			}
		}
		l.Printf("scheduler process complete")
	}()

	return closeSignal, checkSignal, onClosed
}
