package task

import (
	"errors"
	"fmt"
	"time"
)

// Task is a single task struct
type Task struct {
	name          string
	description   string
	completedTime time.Time
	clearedTime   time.Time
}

// New instantiates a new task entity
func New(name string, description string) *Task {
	return &Task{
		name:        name,
		description: description,
	}
}

// NewRaw instantiates a new task entity with all available fields
func NewRaw(name string, description string, complete time.Time, cleared time.Time) *Task {
	return &Task{
		name:          name,
		description:   description,
		completedTime: complete,
		clearedTime:   cleared,
	}
}

// IsValid returns whether a task is valid and can be operated upon
func (t *Task) IsValid() bool {
	return t.clearedTime.IsZero()
}

// Name returns the task namee
func (t *Task) Name() string {
	return t.name
}

// Description returns the task description
func (t *Task) Description() string {
	return t.description
}

// CompletedTime returns the task completed time, zero value if not set
func (t *Task) CompletedTime() time.Time {
	return t.completedTime
}

// ClearedTime returns the task cleared time, zero value if not set
func (t *Task) ClearedTime() time.Time {
	return t.clearedTime
}

// CompleteNow completes a task now
func (t *Task) CompleteNow() (bool, error) {
	if !t.IsValid() {
		return false, errors.New("Task is invalid, cannot be completed")
	}
	if !t.completedTime.IsZero() {
		return false, nil
	}
	t.completedTime = time.Now()
	return true, nil
}

// ClearCompleted clears a completed task now
func (t *Task) ClearCompleted() error {
	if t.completedTime.IsZero() {
		return fmt.Errorf("Incomplete task cannot be cleared: %v", t)
	}
	return t.Clear()
}

// Clear clears a task now if it hasn't previously been cleared
func (t *Task) Clear() error {
	if t.clearedTime.IsZero() {
		t.clearedTime = time.Now()
	}
	return nil
}
