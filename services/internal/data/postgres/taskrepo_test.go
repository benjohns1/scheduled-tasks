// +build integration

package postgres_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestNewTaskRepo(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
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
			}
			gotTasks, err := gotRepo.GetAll()
			if !reflect.DeepEqual(gotTasks, tt.wantTasks) {
				t.Errorf("NewTaskRepo() tasks = %v, want %v", gotTasks, tt.wantTasks)
			}
		})
	}
}

func TestTaskRepo_Get(t *testing.T) {
	now := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	defer clock.Set(prevClock)

	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)
	userRepo, _ := NewUserRepo(conn)
	u := user.New("test user for task Get")
	userRepo.AddExternal(u, "p1", "e1")

	newTask := task.New("t1", "t1desc", u.ID())
	id, err := r.Add(newTask)
	if err != nil {
		t.Fatal(err)
	}

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
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepo_GetAll(t *testing.T) {
	now := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	defer clock.Set(prevClock)

	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)
	userRepo, _ := NewUserRepo(conn)
	u := user.New("test user for task GetAll")
	userRepo.AddExternal(u, "p1", "e1")
	uid := u.ID()

	newTask1 := task.New("", "", uid)
	newTask2 := task.New("", "", uid)
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
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepo_Add(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)
	userRepo, _ := NewUserRepo(conn)
	u := user.New("test user for task Add")
	userRepo.AddExternal(u, "p1", "e1")

	t1 := task.New("", "", user.ID{})
	t2 := task.New("", "", user.New("unknown db user").ID())
	t3 := task.New("t3", "", u.ID())

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
			args:    args{t: t1},
			want:    1,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should return error with unknown user ID",
			r:       r,
			args:    args{t: t2},
			want:    0,
			wantErr: usecase.ErrRecordNotFound,
		},
		{
			name:    "should add 1 task",
			r:       r,
			args:    args{t: t3},
			want:    3,
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Add(tt.args.t)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepo.Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepo_Update(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, _ := NewTaskRepo(conn)
	userRepo, _ := NewUserRepo(conn)
	u := user.New("test user for task Update")
	userRepo.AddExternal(u, "p1", "e1")

	t1 := task.New("t1", "", u.ID())
	id1, _ := r.Add(t1)
	t1.CompleteNow()

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
			args:    args{id: id1, t: t1},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Update(tt.args.id, tt.args.t)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("TaskRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
