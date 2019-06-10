package schedule

// RecurringTask represents a task that recurs on a schedule
type RecurringTask struct {
	name        string
	description string
}

// NewRecurringTask instantiates a new recurring task entity
func NewRecurringTask(name string, description string) RecurringTask {
	return RecurringTask{name, description}
}

// Name returns the task namee
func (rt *RecurringTask) Name() string {
	return rt.name
}

// Description returns the task description
func (rt *RecurringTask) Description() string {
	return rt.description
}

// Equal returns whether 2 recurring tasks are equal
func (rt *RecurringTask) Equal(rtc RecurringTask) bool {
	return rt.name == rtc.name && rt.description == rtc.description
}
