package transient

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestNewTaskRepo(t *testing.T) {
	tests := []struct {
		name string
		want *TaskRepo
	}{
		{
			name: "should return new empty repo",
			want: &TaskRepo{lastID: 0, tasks: make(map[usecase.TaskID]*task.Task)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskRepo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepo_Get(t *testing.T) {
	r := NewTaskRepo()
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
	r := NewTaskRepo()
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
	r := NewTaskRepo()
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
	r := NewTaskRepo()
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
