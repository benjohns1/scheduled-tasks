package schedule

import (
	"fmt"
	"time"
)

// Schedule represents a collection of tasks that recur at some frequency
type Schedule struct {
	frequency Frequency
	paused    bool
	tasks     []RecurringTask
}

// New instantiates a new schedule entity
func New(f Frequency) *Schedule {
	return &Schedule{frequency: f, paused: false, tasks: []RecurringTask{}}
}

// NewRaw creates a new schedule entity from raw data
func NewRaw(frequency Frequency, paused bool, tasks []RecurringTask) *Schedule {
	return &Schedule{frequency, paused, tasks}
}

// Pause pauses a schedule
func (s *Schedule) Pause() {
	s.paused = true
}

// Unpause unpauses a schedule
func (s *Schedule) Unpause() {
	s.paused = false
}

// Paused returns whether schedule is currently paused
func (s *Schedule) Paused() bool {
	return s.paused
}

// Tasks returns the slice of recurring tasks associated with a schedule
func (s *Schedule) Tasks() []RecurringTask {
	return s.tasks
}

// Frequency returns the frequency of this schedule
func (s *Schedule) Frequency() Frequency {
	return s.frequency
}

// AddTask adds a new recurring task if it doesn't already exist
func (s *Schedule) AddTask(rt RecurringTask) error {
	for _, t := range s.tasks {
		if t.Equal(rt) {
			return fmt.Errorf("error adding recurring task: identical task already exists for this schedule")
		}
	}
	s.tasks = append(s.tasks, rt)
	return nil
}

// RemoveTask removes an existing recurring task from the schedule
func (s *Schedule) RemoveTask(rt RecurringTask) error {
	index := -1
	for i, t := range s.tasks {
		if t.Equal(rt) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("error removing task: no matching task found")
	}

	s.tasks = append(s.tasks[:index], s.tasks[index+1:]...)
	return nil
}

// Times gets a list of scheduled times between the start and end times
func (s *Schedule) Times(start time.Time, end time.Time) ([]time.Time, error) {
	return s.frequency.times(start, end)
}
