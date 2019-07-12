package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/pqerr"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// TaskRepo handles persisting task data and maintaining an in-memory cache
type TaskRepo struct {
	db *sql.DB
}

// NewTaskRepo instantiates a new TaskRepo
func NewTaskRepo(conn DBConn) (repo *TaskRepo, err error) {

	if conn.DB == nil {
		return nil, fmt.Errorf("DB connection is nil")
	}

	return &TaskRepo{db: conn.DB}, nil
}

// Get retrieves a task entity, given its persistent ID
func (r *TaskRepo) Get(id usecase.TaskID) (*task.Task, usecase.Error) {

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

	return td.Task, nil
}

// GetForUser retrieves a task entity, given its persistent ID and user ID
func (r *TaskRepo) GetForUser(id usecase.TaskID, uid user.ID) (*task.Task, usecase.Error) {

	// Retrieve from DB
	query := fmt.Sprintf("%s WHERE id = $1 AND created_by = $2", taskSelectClause())
	row := r.db.QueryRow(query, id, uid.String())
	td, err := parseTaskRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.NewError(usecase.ErrRecordNotFound, "no task found with id = %v", id)
		}
		return nil, usecase.NewError(usecase.ErrUnknown, "error parsing task id %d: %v", id, err)
	}

	return td.Task, nil
}

// GetAll retrieves all tasks
func (r *TaskRepo) GetAll() (map[usecase.TaskID]*task.Task, usecase.Error) {
	// Retrieve from DB
	rows, err := r.db.Query(taskSelectClause())
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving all tasks: %v", err)
	}
	defer rows.Close()

	tasks := map[usecase.TaskID]*task.Task{}
	for rows.Next() {
		td, err := parseTaskRow(rows)
		if err != nil {
			return nil, usecase.NewError(usecase.ErrUnknown, "error parsing task row: %v", err)
		}
		tasks[td.TaskID] = td.Task
	}

	return tasks, nil
}

// GetAllForUser retrieves all tasks for a user
func (r *TaskRepo) GetAllForUser(uid user.ID) (map[usecase.TaskID]*task.Task, usecase.Error) {
	q := fmt.Sprintf("%v WHERE created_by = $1", taskSelectClause())

	// Retrieve from DB
	rows, err := r.db.Query(q, uid.String())
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error retrieving all tasks: %v", err)
	}
	defer rows.Close()

	tasks := map[usecase.TaskID]*task.Task{}
	for rows.Next() {
		td, err := parseTaskRow(rows)
		if err != nil {
			return nil, usecase.NewError(usecase.ErrUnknown, "error parsing task row: %v", err)
		}
		tasks[td.TaskID] = td.Task
	}

	return tasks, nil
}

func taskSelectClause() (selectClause string) {
	return "SELECT id, name, description, completed_time, cleared_time, created_time, created_by FROM task"
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
		createdTime   *string
		createdBy     *string
	}
	err = r.Scan(&row.id, &row.name, &row.description, &row.completedTime, &row.clearedTime, &row.createdTime, &row.createdBy)
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
	createdTime, err := time.Parse(dbTimeFormat, *row.createdTime)
	if err != nil {
		createdTime = time.Time{}
	}
	createdBy := user.ID{}
	if row.createdBy != nil {
		createdBy, _ = user.ParseID(*row.createdBy)
	}

	td.Task = task.NewRaw(row.name, row.description, completedTime, clearedTime, createdTime, createdBy)
	td.TaskID = usecase.TaskID(row.id)

	return
}

// Add adds a task to the persisence layer
func (r *TaskRepo) Add(t *task.Task) (usecase.TaskID, usecase.Error) {
	q := "INSERT INTO task (name, description, completed_time, cleared_time, created_time, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	var id usecase.TaskID
	err := r.db.QueryRow(q, t.Name(), t.Description(), t.CompletedTime(), t.ClearedTime(), t.CreatedTime(), t.CreatedBy().StringPtr()).Scan(&id)
	if err != nil {
		if pqerr.Eq(err, pqerr.ForeignKeyViolation) {
			return 0, usecase.NewError(usecase.ErrRecordNotFound, "error inserting new task: %v", err)
		}
		return 0, usecase.NewError(usecase.ErrUnknown, "error inserting new task: %v", err)
	}

	return id, nil
}

// Update updates a task's persistent data to the given entity values
func (r *TaskRepo) Update(id usecase.TaskID, t *task.Task) usecase.Error {
	q := "UPDATE task SET name = $2, description = $3, completed_time = $4, cleared_time = $5, created_time = $6, created_by = $7 WHERE id = $1 RETURNING id"
	rows, err := r.db.Query(q, id, t.Name(), t.Description(), t.CompletedTime(), t.ClearedTime(), t.CreatedTime(), t.CreatedBy().StringPtr())
	if err != nil {
		if pqerr.Eq(err, pqerr.ForeignKeyViolation) {
			return usecase.NewError(usecase.ErrRecordNotFound, "error inserting new task: %v", err)
		}
		return usecase.NewError(usecase.ErrUnknown, "error updating task id %d: %v", id, err)
	}
	defer rows.Close()
	if !rows.Next() {
		return usecase.NewError(usecase.ErrRecordNotFound, "no task found for id = %v", id)
	}

	return nil
}
