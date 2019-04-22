package taskapp

import (
	"database/sql"
	"fmt"
)

// TaskID is the unique id of a task
type TaskID int64

// Task is a single task data struct
type Task struct {
	ID          TaskID
	Name        string
	Description string
	Complete    bool
}

// AddTask persists a new task
func AddTask(db *sql.DB, task *Task) (TaskID, error) {
	stmt, err := db.Prepare("INSERT INTO task (name, description) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("error preparing add task insert statement: %v", err)
	}
	var taskID TaskID
	err = stmt.QueryRow(task.Name, task.Description).Scan(&taskID)
	if err != nil {
		return 0, fmt.Errorf("error inserting new task: %v", err)
	}
	return taskID, nil
}

// CompleteTask completes an existing task
func CompleteTask(db *sql.DB, id TaskID) error {
	stmt, err := db.Prepare("UPDATE task SET complete = TRUE WHERE id = $1")
	if err != nil {
		return fmt.Errorf("error preparing task complete statement: %v", err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error completing task id %d: %v", id, err)
	}
	return nil
}

// ClearCompleted clears completed tasks
func ClearCompleted(db *sql.DB) error {
	stmt, err := db.Prepare("DELETE FROM task WHERE complete = TRUE")
	if err != nil {
		return fmt.Errorf("error preparing task completion statement: %v", err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("error deleting completed tasks: %v", err)
	}
	return nil
}

// ListTasks returns a slice of all tasks
func ListTasks(db *sql.DB) ([]Task, error) {
	return nil, fmt.Errorf("not implemented")
}
