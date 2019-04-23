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
	stmt, err := db.Prepare("UPDATE task SET complete = NOW() WHERE id = $1 RETURNING id;")
	if err != nil {
		return fmt.Errorf("error preparing task complete statement: %v", err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		return fmt.Errorf("error completing task id %d: %v", id, err)
	}
	if !rows.Next() {
		return fmt.Errorf("no tasks completed for id = %v", id)
	}
	return nil
}

// ClearCompleted clears completed tasks
func ClearCompleted(db *sql.DB) (count int, err error) {
	count = 0
	stmt, err := db.Prepare("DELETE FROM task WHERE complete IS NOT NULL RETURNING id;")
	if err != nil {
		err = fmt.Errorf("error preparing task completion statement: %v", err)
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		err = fmt.Errorf("error deleting completed tasks: %v", err)
		return
	}

	for rows.Next() {
		count++
	}
	return
}

// ListTasks returns a slice of all tasks
func ListTasks(db *sql.DB) ([]Task, error) {
	return nil, fmt.Errorf("not implemented")
}
