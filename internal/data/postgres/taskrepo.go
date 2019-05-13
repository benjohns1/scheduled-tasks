package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// TaskRepo handles persisting task data and maintaining an in-memory cache
type TaskRepo struct {
	db    *sql.DB
	tasks map[usecase.TaskID]*task.Task
}

// NewTaskRepo instantiates a new TaskRepo
func NewTaskRepo(conn DBConn) (repo *TaskRepo, err error) {

	if conn.DB == nil {
		return nil, fmt.Errorf("DB connection is nil")
	}

	return &TaskRepo{db: conn.DB, tasks: make(map[usecase.TaskID]*task.Task)}, nil
}

// Get retrieves a task entity, given its persistent ID
func (r *TaskRepo) Get(id usecase.TaskID) (*task.Task, usecase.Error) {

	// Try to retrieve from cache
	t, ok := r.tasks[id]
	if ok {
		return t, nil
	}

	// Retrieve from DB
	query := fmt.Sprintf("%s WHERE id = $1", taskSelectClause())
	row := r.db.QueryRow(query, id)
	td, err := parseTaskRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.NewError(usecase.ErrRecordNotFound, "no task found with id = %v", id)
		}
		return nil, usecase.NewError(usecase.ErrUnknown, "error parsing task id %d: %v", id, err)
	}

	// Add to cache
	r.tasks[td.TaskID] = td.Task

	return td.Task, nil
}

// GetAll retrieves all tasks
func (r *TaskRepo) GetAll() (map[usecase.TaskID]*task.Task, usecase.Error) {
	// Retrieve from DB
	rows, err := r.db.Query(taskSelectClause())
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving all tasks: %v", err)
	}
	for rows.Next() {
		td, err := parseTaskRow(rows)
		if err != nil {
			return nil, usecase.NewError(usecase.ErrUnknown, "error parsing task row: %v", err)
		}

		// Add to cache
		r.tasks[td.TaskID] = td.Task
	}

	return r.tasks, nil
}

func taskSelectClause() (selectClause string) {
	return "SELECT id, name, description, completed_time, cleared_time FROM task"
}

func parseTaskRow(r scannable) (td usecase.TaskData, err error) {

	td = usecase.TaskData{}

	// Scan into row data structure
	var row struct {
		id            int64
		name          string
		description   string
		completedTime *string
		clearedTime   *string
	}
	err = r.Scan(&row.id, &row.name, &row.description, &row.completedTime, &row.clearedTime)
	if err != nil {
		return
	}

	// Map values
	completedTime, err := time.Parse(dbTimeFormat, *row.completedTime)
	if err != nil {
		completedTime = time.Time{}
	}
	clearedTime, err := time.Parse(dbTimeFormat, *row.clearedTime)
	if err != nil {
		clearedTime = time.Time{}
	}

	td.Task = task.NewRaw(row.name, row.description, completedTime, clearedTime)
	td.TaskID = usecase.TaskID(row.id)

	return
}

// Add adds a task to the persisence layer
func (r *TaskRepo) Add(t *task.Task) (usecase.TaskID, usecase.Error) {
	q := "INSERT INTO task (name, description, completed_time, cleared_time) VALUES ($1, $2, $3, $4) RETURNING id"
	var id usecase.TaskID
	err := r.db.QueryRow(q, t.Name(), t.Description(), t.CompletedTime(), t.ClearedTime()).Scan(&id)
	if err != nil {
		return 0, usecase.NewError(usecase.ErrUnknown, "error inserting new task: %v", err)
	}

	r.tasks[id] = t

	return id, nil
}

// Update updates a task's persistent data to the given entity values
func (r *TaskRepo) Update(id usecase.TaskID, t *task.Task) usecase.Error {
	q := "UPDATE task SET name = $1, description = $2, completed_time = $3, cleared_time = $4 WHERE id = $5 RETURNING id"
	rows, err := r.db.Query(q, t.Name(), t.Description(), t.CompletedTime(), t.ClearedTime(), id)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error updating task id %d: %v", id, err)
	}
	if !rows.Next() {
		return usecase.NewError(usecase.ErrRecordNotFound, "no task found for id = %v", id)
	}

	r.tasks[id] = t

	return nil
}
