package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// TaskRepo handles persisting task data and maintaining an in-memory cache
type TaskRepo struct {
	l     Logger
	db    *sql.DB
	tasks map[usecase.TaskID]*core.Task
	Close func()
}

const dbTimeFormat = time.RFC3339Nano

// Logger interface needed for log messages
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// NewTaskRepo instantiates a new TaskRepo
func NewTaskRepo(l Logger, conn DBConn) (repo *TaskRepo, err error) {

	close := func() {}

	// Connect to DB
	l.Printf("connecting to db %s as %s...", conn.Name, conn.User)
	db, err := connect(l, conn)
	if db != nil {
		close = func() {
			db.Close()
		}
	}
	if err != nil {
		close()
		err = fmt.Errorf("error opening db: %v", err)
		return
	}

	// Perform DB setup if needed
	didSetup, err := setup(db)
	if err != nil {
		close()
		err = fmt.Errorf("error setting up db: %v", err)
		return
	}
	if didSetup {
		l.Print("first-time DB setup complete")
	}

	repo = &TaskRepo{l: l, db: db, tasks: make(map[usecase.TaskID]*core.Task), Close: close}

	return
}

// WipeAndReset completely destroys all data in persistence and cache
func (r *TaskRepo) WipeAndReset() usecase.Error {

	// Drop db table
	_, err := r.db.Exec("DROP TABLE task")
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error dropping task table: %v", err)
	}

	// Reset db table
	_, err = setup(r.db)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error resetting task table: %v", err)
	}

	// Destroy/reset cache
	r.tasks = make(map[usecase.TaskID]*core.Task)

	return nil
}

// Get retrieves a task entity, given its persistent ID
func (r *TaskRepo) Get(id usecase.TaskID) (*core.Task, usecase.Error) {

	// Try to retrieve from cache
	t, ok := r.tasks[id]
	if ok {
		return t, nil
	}

	// Retrieve from DB
	row := r.db.QueryRow("SELECT id, name, description, completed_time, cleared_time FROM task WHERE id = $1", id)
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
func (r *TaskRepo) GetAll() (map[usecase.TaskID]*core.Task, usecase.Error) {
	// Retrieve from DB
	rows, err := r.db.Query("SELECT id, name, description, completed_time, cleared_time FROM task")
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

type scannable interface {
	Scan(dest ...interface{}) error
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

	td.Task = core.NewTaskFull(row.name, row.description, completedTime, clearedTime)
	td.TaskID = usecase.TaskID(row.id)

	return
}

// Add adds a task to the persisence layer
func (r *TaskRepo) Add(t *core.Task) (usecase.TaskID, usecase.Error) {
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
func (r *TaskRepo) Update(id usecase.TaskID, t *core.Task) usecase.Error {
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
