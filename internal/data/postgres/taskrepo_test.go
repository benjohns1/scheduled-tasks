// +build integration

package postgres

import (
	"os"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/benjohns1/scheduled-tasks/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
	"github.com/joho/godotenv"
)

type LoggerMock struct{}

func (l *LoggerMock) Print(v ...interface{})                 {}
func (l *LoggerMock) Printf(format string, v ...interface{}) {}
func (l *LoggerMock) Println(v ...interface{})               {}

func mockDBConn() (DBConn, error) {
	l := &LoggerMock{}

	// Load environment vars
	if err := godotenv.Load("../../../.env"); err != nil {
		return DBConn{}, fmt.Errorf("could not load .env file: %v", err)
	}

	// Load DB connection info
	dbconn := NewDBConn(l)
	testPort, err := strconv.Atoi(os.Getenv("POSTGRES_TEST_PORT"))
	if err != nil || testPort == 0 {
		dbconn.Close()
		return dbconn, fmt.Errorf("POSTGRES_TEST_PORT must be set to run postgres DB integreation tests %v", err)
	}
	dbconn.Port = testPort
	dbconn.MaxRetryAttempts = 1
	dbconn.RetrySleepSeconds = 0

	// Connect to DB, destroy any existing tables, setup tables again
	if err := dbconn.Connect(); err != nil {
		dbconn.Close()
		return dbconn, fmt.Errorf("could not connect to test DB: %v", err)
	}
	_, detroyErr := dbconn.destroy()
	didSetup, err := dbconn.Setup()
	if err != nil {
		dbconn.Close()
		return dbconn, err
	}
	if !didSetup {
		dbconn.Close()
		return dbconn, fmt.Errorf("could not setup fresh DB tables, test tables may not have been not properly destroyed: %v", detroyErr)
	}
	return dbconn, nil
}

func TestNewTaskRepo(t *testing.T) {
	conn, err := mockDBConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	type args struct {
		conn DBConn
	}
	tests := []struct {
		name      string
		args      args
		wantTasks map[usecase.TaskID]*task.Task
		wantErr   bool
	}{
		{
			name:      "should return new empty repo",
			args:      args{conn},
			wantTasks: map[usecase.TaskID]*task.Task{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, err := NewTaskRepo(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTaskRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRepo.tasks, tt.wantTasks) {
				t.Errorf("NewTaskRepo() tasks = %v, want %v", gotRepo.tasks, tt.wantTasks)
			}
		})
	}
}

func TestTaskRepo_Get(t *testing.T) {
	conn, err := mockDBConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)

	newTask := task.New("", "")
	id, _ := r.Add(newTask)

	type args struct {
		id usecase.TaskID
	}
	tests := []struct {
		name    string
		r       *TaskRepo
		args    args
		want    *task.Task
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 1 empty task",
			r:       r,
			args:    args{id: id},
			want:    newTask,
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Get(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.Get() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTaskRepo_GetAll(t *testing.T) {
	conn, err := mockDBConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)

	newTask1 := task.New("", "")
	newTask2 := task.New("", "")
	id1, _ := r.Add(newTask1)
	id2, _ := r.Add(newTask2)

	tests := []struct {
		name    string
		r       *TaskRepo
		want    map[usecase.TaskID]*task.Task
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 2 empty tasks",
			r:       r,
			want:    map[usecase.TaskID]*task.Task{id1: newTask1, id2: newTask2},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAll()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.GetAll() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTaskRepo_Add(t *testing.T) {
	conn, err := mockDBConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)

	newTask := task.New("", "")

	type args struct {
		t *task.Task
	}
	tests := []struct {
		name    string
		r       *TaskRepo
		args    args
		want    usecase.TaskID
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should add 1 empty task",
			r:       r,
			args:    args{t: newTask},
			want:    1,
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Add(tt.args.t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.Add() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTaskRepo_Update(t *testing.T) {
	conn, err := mockDBConn()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)

	newTask := task.New("", "")
	id1, _ := r.Add(newTask)
	newTask.CompleteNow()

	type args struct {
		id usecase.TaskID
		t  *task.Task
	}
	tests := []struct {
		name    string
		r       *TaskRepo
		args    args
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should successfully update task",
			r:       r,
			args:    args{id: id1, t: newTask},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Update(tt.args.id, tt.args.t)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
