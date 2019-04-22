package taskapp

import (
	"database/sql"
	"fmt"
)

// TaskKey is the unique key of a task
type TaskKey int

// Task is a single task data struct
type Task struct {
	Key         TaskKey
	Name        string
	Description string
}

// AddTask persists a new task
func AddTask(db *sql.DB, task *Task) (TaskKey, error) {
	stmt, err := db.Prepare("INSERT INTO task (name, description) VALUES ($1, $2) RETURNING key;")
	if err != nil {
		return 0, fmt.Errorf("error preparing add task insert statement: %v", err)
	}
	var taskKey TaskKey
	stmt.Exec()
	err = stmt.QueryRow(task.Name, task.Description).Scan(&taskKey)
	if err != nil {
		return 0, fmt.Errorf("error inserting new task: %v", err)
	}
	return taskKey, nil
}

// CompleteTask completes an existing task
func CompleteTask(db *sql.DB, key TaskKey) error {
	return fmt.Errorf("not implemented")
}

// ClearCompleted clears completed tasks
func ClearCompleted(db *sql.DB) error {
	return fmt.Errorf("not implemented")
}

// ListTasks returns a slice of all tasks
func ListTasks(db *sql.DB) ([]Task, error) {
	return nil, fmt.Errorf("not implemented")
}
