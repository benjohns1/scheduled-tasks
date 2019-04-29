package usecase_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	data "github.com/benjohns1/scheduled-tasks/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/internal/usecase"
)

func TestAddTask(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}

	type args struct {
		r TaskRepo
		t *core.Task
	}
	tests := []struct {
		name    string
		args    args
		want    *TaskData
		wantErr bool
	}{
		{
			name:    "add empty task",
			args:    args{r: taskRepo, t: core.NewTask("", "")},
			want:    &TaskData{Task: core.NewTask("", "")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddTask(tt.args.r, tt.args.t)
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
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	taskID, err1 := taskRepo.Add(core.NewTask("task1", ""))
	completedTaskID, err2 := taskRepo.Add(core.NewTaskFull("task2", "", time.Now(), time.Time{}))
	clearedTaskID, err3 := taskRepo.Add(core.NewTaskFull("task3", "", time.Now(), time.Now()))
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("error setting up task repo: %v, %v, %v", err1, err2, err3)
	}

	type args struct {
		r  TaskRepo
		id TaskID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "task should be completed",
			args:    args{r: taskRepo, id: taskID},
			want:    true,
			wantErr: false,
		},
		{
			name:    "completed task should not be completed again",
			args:    args{r: taskRepo, id: completedTaskID},
			want:    false,
			wantErr: false,
		},
		{
			name:    "completing a cleared task should return error",
			args:    args{r: taskRepo, id: clearedTaskID},
			want:    false,
			wantErr: true,
		},
		{
			name:    "completing an invalid task ID should return error",
			args:    args{r: taskRepo, id: 99999},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompleteTask(tt.args.r, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CompleteTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearTask(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	taskID, err1 := taskRepo.Add(core.NewTask("task1", ""))
	completedTaskID, err2 := taskRepo.Add(core.NewTaskFull("task2", "", time.Now(), time.Time{}))
	clearedTaskID, err3 := taskRepo.Add(core.NewTaskFull("task3", "", time.Now(), time.Now()))
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("error setting up task repo: %v, %v, %v", err1, err2, err3)
	}

	type args struct {
		r  TaskRepo
		id TaskID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "task should be cleared",
			args:    args{r: taskRepo, id: taskID},
			want:    true,
			wantErr: false,
		},
		{
			name:    "completed task should be cleared",
			args:    args{r: taskRepo, id: completedTaskID},
			want:    true,
			wantErr: false,
		},
		{
			name:    "previously cleared task should return false",
			args:    args{r: taskRepo, id: clearedTaskID},
			want:    false,
			wantErr: false,
		},
		{
			name:    "clearing a non-existent task ID should return error",
			args:    args{r: taskRepo, id: 9999},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClearTask(tt.args.r, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ClearTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearCompletedTasks(t *testing.T) {
	emptyRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}

	singleCompletedTaskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	_, err1 := singleCompletedTaskRepo.Add(core.NewTask("task1", ""))
	_, err2 := singleCompletedTaskRepo.Add(core.NewTaskFull("task2", "", time.Now(), time.Time{}))
	_, err3 := singleCompletedTaskRepo.Add(core.NewTaskFull("task3", "", time.Now(), time.Now()))
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("error setting up task repo: %v, %v, %v", err1, err2, err3)
	}

	thousandCompletedTasksRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	for i := 0; i < 1000; i++ {
		_, err := thousandCompletedTasksRepo.Add(core.NewTaskFull("", "", time.Now(), time.Time{}))
		if err != nil {
			t.Errorf("error setting up task repo: %v", err)
		}
	}

	type args struct {
		r TaskRepo
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "clearing empty repo should return 0 count",
			args:      args{r: emptyRepo},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "clearing repo should return 1 count",
			args:      args{r: singleCompletedTaskRepo},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "clearing repo should return 1000 count",
			args:      args{r: thousandCompletedTasksRepo},
			wantCount: 1000,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := ClearCompletedTasks(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClearCompletedTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ClearCompletedTasks() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	task1 := core.NewTask("task1", "")
	id1, err1 := taskRepo.Add(task1)
	task2 := core.NewTaskFull("task2", "", time.Now(), time.Time{})
	id2, err2 := taskRepo.Add(task2)
	_, err3 := taskRepo.Add(core.NewTaskFull("task3", "", time.Now(), time.Now()))
	if err1 != nil || err2 != nil || err3 != nil {
		t.Errorf("error setting up task repo: %v, %v, %v", err1, err2, err3)
	}

	type args struct {
		r TaskRepo
	}
	tests := []struct {
		name    string
		args    args
		want    map[TaskID]*core.Task
		wantErr bool
	}{
		{
			name: "task list should return map of 2 uncleared tasks",
			args: args{r: taskRepo},
			want: map[TaskID]*core.Task{
				id1: task1,
				id2: task2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTasks(tt.args.r)
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