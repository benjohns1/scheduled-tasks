package usecase

import (
	"fmt"

	"github.com/benjohns1/scheduled-tasks/internal/core"
)

// TaskID is the persistent ID of the task
type TaskID int64

// TaskData contains appication-level task info
type TaskData struct {
	TaskID TaskID
	Task   *core.Task
}

// TaskRepo defines the task repository interface required by use cases
type TaskRepo interface {
	Get(TaskID) (*core.Task, error)
	GetAll() (map[TaskID]*core.Task, error)
	Add(*core.Task) (TaskID, error)
	Update(TaskID, *core.Task) error
	WipeAndReset() error
}

// AddTask creates and adds a new task to the list
func AddTask(r TaskRepo, name string, desc string) (*TaskData, error) {
	task := core.NewTask(name, desc)
	id, err := r.Add(task)
	if err != nil {
		return nil, fmt.Errorf("error adding task: %v", err)
	}
	taskData := &TaskData{TaskID: id, Task: task}
	return taskData, nil
}

// CompleteTask completes an existing task
func CompleteTask(r TaskRepo, id TaskID) (bool, error) {
	t, err := r.Get(id)
	if err != nil {
		return false, fmt.Errorf("error retrieving task id %d: %v", id, err)
	}

	ok, err := t.CompleteNow()
	if err != nil {
		return false, fmt.Errorf("error completing task id %d: %v", id, err)
	}
	if !ok {
		return false, nil
	}

	err = r.Update(id, t)
	if err != nil {
		return false, fmt.Errorf("error updating task id %d: %v", id, err)
	}
	return true, nil
}

// ClearTask clears (removes) a single task, regardless of whether it has been completed
func ClearTask(r TaskRepo, id TaskID) (bool, error) {
	t, err := r.Get(id)
	if err != nil {
		return false, fmt.Errorf("error retrieving task id %d: %v", id, err)
	}

	err = t.Clear()
	if err != nil {
		return false, fmt.Errorf("error clearing task id %d: %v", id, err)
	}

	err = r.Update(id, t)
	if err != nil {
		return false, fmt.Errorf("error updating task id %d: %v", id, err)
	}
	return true, nil
}

// ClearCompletedTasks clears all completed tasks
func ClearCompletedTasks(r TaskRepo) (count int, err error) {
	count = 0

	ts, err := r.GetAll()
	if err != nil {
		err = fmt.Errorf("error retrieving tasks to clear: %v", err)
		return
	}

	for id, t := range ts {
		if t.CompletedTime().IsZero() || !t.ClearedTime().IsZero() {
			continue
		}
		err = t.ClearCompleted()
		if err != nil {
			return
		}
		err = r.Update(id, t)
		if err != nil {
			return
		}
		count++
	}

	return
}

// ListTasks returns all tasks that haven't been cleared
func ListTasks(r TaskRepo) (map[TaskID]*core.Task, error) {

	all, err := r.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error retrieving tasks: %v", err)
	}

	list := make(map[TaskID]*core.Task)
	for id, t := range all {
		if !t.ClearedTime().IsZero() {
			continue
		}
		list[id] = t
	}

	return list, nil
}
