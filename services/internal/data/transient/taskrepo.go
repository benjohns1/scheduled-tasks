package transient

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// TaskRepo maintains an in-memory cache of tasks
type TaskRepo struct {
	lastID int
	tasks  map[usecase.TaskID]*task.Task
}

// NewTaskRepo instantiates a new TaskRepo
func NewTaskRepo() *TaskRepo {
	return &TaskRepo{tasks: make(map[usecase.TaskID]*task.Task)}
}

// Get retrieves a task entity, given its persistent ID
func (r *TaskRepo) Get(id usecase.TaskID) (*task.Task, usecase.Error) {

	// Try to retrieve from cache
	t, ok := r.tasks[id]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no task with ID: %v", id)
	}
	return t, nil
}

// GetAll retrieves all tasks
func (r *TaskRepo) GetAll() (map[usecase.TaskID]*task.Task, usecase.Error) {

	return r.tasks, nil
}

// Add adds a task to the persisence layer
func (r *TaskRepo) Add(t *task.Task) (usecase.TaskID, usecase.Error) {
	r.lastID++
	id := usecase.TaskID(r.lastID)
	r.tasks[id] = t

	return id, nil
}

// Update updates a task's persistent data to the given entity values
func (r *TaskRepo) Update(id usecase.TaskID, t *task.Task) usecase.Error {

	_, ok := r.tasks[id]
	if !ok {
		return usecase.NewError(usecase.ErrRecordNotFound, "no task with ID %v", id)
	}

	r.tasks[id] = t

	return nil
}
