package scheduler

import (
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

const offset = 3 * time.Second

type nextRun struct {
	next time.Time
	halt bool
}

// Run starts the scheduler process
func Run(l Logger, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo) {
	l.Printf("scheduler process starting")

	for {
		nextRun := checkSchedules(l)
		if nextRun.halt {
			break
		}

		// Sleep until next scheduled time + offset
		wait := nextRun.next.Sub(time.Now())
		time.Sleep(wait + offset)
	}

	l.Printf("scheduler process complete")
}

func checkSchedules(l Logger) nextRun {
	l.Printf("checking schedules")
	now := time.Now()
	l.Printf("Now: %v", now)

	// TODO: Check all schedules to see which have at least one recurrence between now and the last time they recurred
	// TODO: Kick-off async processing for each of these
	// TODO: Of the remaining schedules, get the time the soonest one is scheduled to recur

	return nextRun{
		next: time.Now().Add(1 * time.Second),
		halt: false,
	}
}
