package usecase

import (
	"fmt"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core/task"
)

// CheckSchedules checks all schedules, determines all recurrences that have occurred, and when the next run is needed
func CheckSchedules(c Clock, taskRepo TaskRepo, scheduleRepo ScheduleRepo) (time.Time, error) {
	//now := timeNow()

	// Check all schedules to see which have at least one recurrence between now and the last time they recurred
	schedules, err := scheduleRepo.GetAll()
	if err != nil {
		return time.Time{}, err
	}

	now := c.Now()
	var next time.Time
	for id, sched := range schedules {
		// TODO: persistently store when last successful run occured for each schedule and use as 'start' param (new schedule_last_run table or reuse schedule table?)
		times, err := sched.Times(now.Add(-15*time.Minute), now)
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
