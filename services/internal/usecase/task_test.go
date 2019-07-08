package usecase_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	data "github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestGetTask(t *testing.T) {
	now := clock.Now()
	taskRepo := data.NewTaskRepo()
	task1 := task.New("task1", "", user.ID{})
	taskID, _ := taskRepo.Add(task1)
	task2 := task.NewRaw("task2", "", now, time.Time{}, now, user.ID{})
	completedTaskID, _ := taskRepo.Add(task2)
	clearedTaskID, _ := taskRepo.Add(task.NewRaw("task3", "", now, now, now, user.ID{}))

	type args struct {
		r  TaskRepo
		id TaskID
	}
	tests := []struct {
		name    string
		args    args
		want    *TaskData
		wantErr ErrorCode
	}{
		{
			name:    "task should be retrieved",
			args:    args{r: taskRepo, id: taskID},
			want:    &TaskData{TaskID: taskID, Task: task1},
			wantErr: ErrNone,
		},
		{
			name:    "completed task should be retrieved",
			args:    args{r: taskRepo, id: completedTaskID},
			want:    &TaskData{TaskID: completedTaskID, Task: task2},
			wantErr: ErrNone,
		},
		{
			name:    "cleared task should return an ErrRecordNotFound",
			args:    args{r: taskRepo, id: clearedTaskID},
			want:    nil,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "non-existent task ID should return an ErrRecordNotFound",
			args:    args{r: taskRepo, id: 99999},
			want:    nil,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTask(tt.args.r, tt.args.id)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	taskRepo := data.NewTaskRepo()
	emptyTask := task.New("", "", user.ID{})
	basicTask := task.New("task with data", "task description", user.ID{})

	type args struct {
		r TaskRepo
		t *task.Task
	}
	tests := []struct {
		name    string
		args    args
		want    *task.Task
		wantErr bool
	}{
		{
			name:    "add empty task should be valid",
			args:    args{r: taskRepo, t: emptyTask},
			want:    emptyTask,
			wantErr: false,
		},
		{
			name:    "add task with basic info should be valid",
			args:    args{r: taskRepo, t: basicTask},
			want:    basicTask,
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
			if !reflect.DeepEqual(got.Task, tt.want) {
				t.Errorf("AddTask() = %v, want %v", got.Task, tt.want)
			}
		})
	}
}

func TestCompleteTask(t *testing.T) {
	now := clock.Now()
	taskRepo := data.NewTaskRepo()
	taskID, _ := taskRepo.Add(task.New("task1", "", user.ID{}))
	completedTaskID, _ := taskRepo.Add(task.NewRaw("task2", "", now, time.Time{}, now, user.ID{}))
	clearedTaskID, _ := taskRepo.Add(task.NewRaw("task3", "", now, now, now, user.ID{}))

	type args struct {
		r  TaskRepo
		id TaskID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr ErrorCode
	}{
		{
			name:    "task should be completed",
			args:    args{r: taskRepo, id: taskID},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "completed task should not be completed again",
			args:    args{r: taskRepo, id: completedTaskID},
			want:    false,
			wantErr: ErrNone,
		},
		{
			name:    "completing a cleared task should return an ErrRecordNotFound",
			args:    args{r: taskRepo, id: clearedTaskID},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "completing a non-existent task ID should return an ErrRecordNotFound",
			args:    args{r: taskRepo, id: 99999},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompleteTask(tt.args.r, tt.args.id)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
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
	now := clock.Now()
	taskRepo := data.NewTaskRepo()
	taskID, _ := taskRepo.Add(task.New("task1", "", user.ID{}))
	completedTaskID, _ := taskRepo.Add(task.NewRaw("task2", "", now, time.Time{}, now, user.ID{}))
	clearedTaskID, _ := taskRepo.Add(task.NewRaw("task3", "", now, now, now, user.ID{}))

	type args struct {
		r  TaskRepo
		id TaskID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr ErrorCode
	}{
		{
			name:    "task should be cleared",
			args:    args{r: taskRepo, id: taskID},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "completed task should be cleared",
			args:    args{r: taskRepo, id: completedTaskID},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "previously cleared task should return ErrRecordNotFound",
			args:    args{r: taskRepo, id: clearedTaskID},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "clearing a non-existent task ID should return ErrRecordNotFound",
			args:    args{r: taskRepo, id: 9999},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClearTask(tt.args.r, tt.args.id)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
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
	now := clock.Now()
	emptyRepo := data.NewTaskRepo()

	singleCompletedTaskRepo := data.NewTaskRepo()
	singleCompletedTaskRepo.Add(task.New("task1", "", user.ID{}))
	singleCompletedTaskRepo.Add(task.NewRaw("task2", "", now, time.Time{}, now, user.ID{}))
	singleCompletedTaskRepo.Add(task.NewRaw("task3", "", now, now, now, user.ID{}))

	thousandCompletedTasksRepo := data.NewTaskRepo()
	for i := 0; i < 1000; i++ {
		_, err := thousandCompletedTasksRepo.Add(task.NewRaw("", "", now, time.Time{}, now, user.ID{}))
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
	now := clock.Now()
	taskRepo := data.NewTaskRepo()
	task1 := task.New("task1", "", user.ID{})
	id1, _ := taskRepo.Add(task1)
	task2 := task.NewRaw("task2", "", now, time.Time{}, now, user.ID{})
	id2, _ := taskRepo.Add(task2)
	taskRepo.Add(task.NewRaw("task3", "", now, now, now, user.ID{}))

	type args struct {
		r TaskRepo
	}
	tests := []struct {
		name    string
		args    args
		want    map[TaskID]*task.Task
		wantErr bool
	}{
		{
			name: "task list should return map of 2 uncleared tasks",
			args: args{r: taskRepo},
			want: map[TaskID]*task.Task{
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
