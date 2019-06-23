package usecase

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
)

// TaskID is the persistent ID of the task
type TaskID int64

// TaskData contains application-level task info
type TaskData struct {
	TaskID TaskID
	Task   *task.Task
}

// TaskRepo defines the task repository interface required by use cases
type TaskRepo interface {
	Get(TaskID) (*task.Task, Error)
	GetAll() (map[TaskID]*task.Task, Error)
	Add(*task.Task) (TaskID, Error)
	Update(TaskID, *task.Task) Error
}

// GetTask gets a single task
func GetTask(r TaskRepo, id TaskID) (*TaskData, Error) {
	t, ucerr := r.Get(id)
	if ucerr != nil {
		return nil, ucerr.Prefix("error retrieving task id %d", id)
	}
	if !t.IsValid() {
		return nil, NewError(ErrRecordNotFound, "task id %d not found", id)
	}
	return &TaskData{TaskID: id, Task: t}, nil
}

// AddTask creates and adds a new task to the list
func AddTask(r TaskRepo, t *task.Task) (*TaskData, Error) {
	id, err := r.Add(t)
	if err != nil {
		return nil, NewError(ErrUnknown, "error adding task: %v", err)
	}
	taskData := &TaskData{TaskID: id, Task: t}
	return taskData, nil
}

// CompleteTask completes an existing task
func CompleteTask(r TaskRepo, id TaskID) (bool, Error) {
	t, ucerr := r.Get(id)
	if ucerr != nil {
		return false, ucerr.Prefix("error retrieving task id %d", id)
	}

	if !t.IsValid() {
		return false, NewError(ErrRecordNotFound, "task id %d not found", id)
	}

	ok, err := t.CompleteNow()
	if err != nil {
		return false, NewError(ErrUnknown, "error completing task id %d: %v", id, err)
	}
	if !ok {
		return false, nil
	}

	ucerr = r.Update(id, t)
	if ucerr != nil {
		return false, ucerr.Prefix("error updating task id %d", id)
	}
	return true, nil
}

// ClearTask clears (removes) a single task, regardless of whether it has been completed
func ClearTask(r TaskRepo, id TaskID) (bool, Error) {
	t, ucerr := r.Get(id)
	if ucerr != nil {
		if ucerr.Code() == ErrRecordNotFound {
			return false, ucerr.Prefix("task id %d not found", id)
		}
		return false, ucerr.Prefix("error retrieving task id %d", id)
	}

	if !t.IsValid() {
		return false, NewError(ErrRecordNotFound, "task id %d not found", id)
	}

	err := t.Clear()
	if err != nil {
		return false, NewError(ErrUnknown, "error clearing task id %d: %v", id, err)
	}

	ucerr = r.Update(id, t)
	if ucerr != nil {
		return false, ucerr.Prefix("error updating task id %d", id)
	}
	return true, nil
}

// ClearCompletedTasks clears all completed tasks, returning the number completed and an error
func ClearCompletedTasks(r TaskRepo) (int, Error) {
	ts, ucerr := r.GetAll()
	if ucerr != nil {
		return 0, ucerr.Prefix("error retrieving tasks to clear")
	}

	count := 0
	for id, t := range ts {
		if t.CompletedTime().IsZero() || !t.IsValid() {
			continue
		}
		err := t.ClearCompleted()
		if err != nil {
			return count, NewError(ErrUnknown, "error clearing completed tasks: %v", err)
		}
		ucerr = r.Update(id, t)
		if ucerr != nil {
			return count, ucerr
		}
		count++
	}

	return count, nil
}

// ListTasks returns all valid (uncleared) tasks
func ListTasks(r TaskRepo) (map[TaskID]*task.Task, Error) {
	all, ucerr := r.GetAll()
	if ucerr != nil {
		return nil, ucerr.Prefix("error retrieving tasks")
	}

	list := make(map[TaskID]*task.Task)
	for id, t := range all {
		if !t.IsValid() {
			continue
		}
		list[id] = t
	}

	return list, nil
}
