// +build integration

package taskapp

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	persistence "github.com/benjohns1/scheduled-tasks/internal/pkg/persistence/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func connect() (db *sql.DB, teardown func() (sql.Result, error), err error) {
	// Load environment vars
	godotenv.Load("../../../.env")

	// Load DB connection info
	connInfo, err := persistence.LoadConnInfo()
	if err != nil {
		err = fmt.Errorf("error loading db connection details: %v", err)
	}

	// Connect to DB
	db, err = persistence.Connect(connInfo)
	if err != nil {
		err = fmt.Errorf("error opening db: %v", err)
	}

	// Perform DB setup if needed
	teardown, err = persistence.TestSetup(db)
	if err != nil {
		err = fmt.Errorf("error setting up db: %v", err)
	}
	return
}

func TestAddTask(t *testing.T) {
	tests := []struct {
		task    *Task
		want    TaskID
		wantErr bool
	}{
		{task: &Task{Name: "task name", Description: "task description"}, want: 1, wantErr: false},
		{task: &Task{Name: "asdf", Description: "1234"}, want: 2, wantErr: false},
	}

	// Get DB connection
	db, teardown, err := connect()
	if db != nil {
		defer db.Close()
		defer teardown()
	}
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
		return
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.task.Name, func(t *testing.T) {

			got, err := AddTask(db, tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompleteTask(t *testing.T) {
	type args struct {
		db *sql.DB
		id TaskID
	}

	// Get DB connection
	db, teardown, err := connect()
	if db != nil {
		defer db.Close()
		defer teardown()
	}
	if err != nil {
		t.Errorf("error connecting to db: %v", err)
		return
	}

	// Seed with 3 tasks
	AddTask(db, &Task{Name: "first task", Description: ""})
	AddTask(db, &Task{Name: "second task", Description: ""})
	AddTask(db, &Task{Name: "third task", Description: ""})

	// Setup tests
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "first", args: args{db: db, id: TaskID(1)}, wantErr: false},
		{name: "second", args: args{db: db, id: TaskID(2)}, wantErr: false},
		{name: "third", args: args{db: db, id: TaskID(3)}, wantErr: false},
		{name: "fourth", args: args{db: db, id: TaskID(4)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CompleteTask(tt.args.db, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CompleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClearCompleted(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := ClearCompleted(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearCompleted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ClearCompleted() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	type args struct {
		db *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []Task
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTasks(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
