package usecase

import (
	"fmt"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
)

// CheckSchedules checks all schedules, determines all recurrences that have occurred, and when the next run is needed
func CheckSchedules(taskRepo TaskRepo, scheduleRepo ScheduleRepo) (time.Time, error) {
	// Check all unpaused schedules
	schedules, err := scheduleRepo.GetAllUnpaused()
	if err != nil {
		return time.Time{}, err
	}

	now := clock.Now()
	var next time.Time
	for id, sched := range schedules {
		// If the schedule has previously been checked, create tasks for any recurrences
		if !sched.LastChecked().IsZero() {

			times, err := sched.Times(sched.LastChecked(), now)
			if err != nil {
				return time.Time{}, fmt.Errorf("error retrieving times from schedule id %v: %v", id, err)
			}

			// Create tasks for all scheduled recurrences
			for _, rt := range sched.Tasks() {
				for i := 0; i < len(times); i++ {
					t := task.New(rt.Name(), rt.Description())
					_, err := taskRepo.Add(t)
					if err != nil {
						return time.Time{}, fmt.Errorf("error adding task to repo: %v", err)
					}
				}
			}
		}

		// Set the last checked time for the schedule
		sched.Check(now)
		scheduleRepo.Update(id, sched)

		// Get the next runtime and store the nearest upcoming time as the next time to run scheduler
		n, err := sched.NextTime(now)
		if err != nil {
			return time.Time{}, fmt.Errorf("error getting next schedule time for id %v: %v", id, err)
		}
		if next.IsZero() || n.Before(next) {
			next = n
		}
	}

	return next, nil
}
