package schedule

import (
	"time"
)

// Schedule represents a collection of tasks that recur at some frequency
type Schedule struct {
	frequency *Frequency
	paused    bool
	tasks     []RecurringTask
}

// New instantiates a new schedule entity
func New(f *Frequency) *Schedule {
	return &Schedule{frequency: f, paused: false, tasks: []RecurringTask{}}
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

// Times gets a list of scheduled times between the start and end times
func (s *Schedule) Times(start time.Time, end time.Time) ([]time.Time, error) {
	return s.frequency.times(start, end)
}
