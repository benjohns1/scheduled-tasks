package transient

import (
	"github.com/benjohns1/scheduled-tasks/internal/core"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// TaskRepo maintains an in-memory cache of tasks
type TaskRepo struct {
	lastID int
	tasks  map[usecase.TaskID]*core.Task
}

// NewTaskRepo instantiates a new TaskRepo
func NewTaskRepo() (repo *TaskRepo, err error) {

	err = nil
	repo = &TaskRepo{tasks: make(map[usecase.TaskID]*core.Task)}

	return
}

// WipeAndReset completely destroys all data in persistence and cache
func (r *TaskRepo) WipeAndReset() usecase.Error {

	// Destroy/reset cache
	r.tasks = make(map[usecase.TaskID]*core.Task)

	return nil
}

// Get retrieves a task entity, given its persistent ID
func (r *TaskRepo) Get(id usecase.TaskID) (*core.Task, usecase.Error) {

	// Try to retrieve from cache
	t, ok := r.tasks[id]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no task with ID: %v", id)
	}
	return t, nil
}

// GetAll retrieves all tasks
func (r *TaskRepo) GetAll() (map[usecase.TaskID]*core.Task, usecase.Error) {

	return r.tasks, nil
}

// Add adds a task to the persisence layer
func (r *TaskRepo) Add(t *core.Task) (usecase.TaskID, usecase.Error) {
	id := usecase.TaskID(r.lastID)
	r.lastID++

	r.tasks[id] = t

	return id, nil
}

// Update updates a task's persistent data to the given entity values
func (r *TaskRepo) Update(id usecase.TaskID, t *core.Task) usecase.Error {

	_, ok := r.tasks[id]
	if !ok {
		return usecase.NewError(usecase.ErrRecordNotFound, "no task with ID %v", id)
	}

	r.tasks[id] = t

	return nil
}
