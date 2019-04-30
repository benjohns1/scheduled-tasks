package usecase_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	data "github.com/benjohns1/scheduled-tasks/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/internal/usecase"
)

func TestGetTask(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}
	task1 := core.NewTask("task1", "")
	taskID, err1 := taskRepo.Add(task1)
	task2 := core.NewTaskFull("task2", "", time.Now(), time.Time{})
	completedTaskID, err2 := taskRepo.Add(task2)
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
				t.Errorf("CompleteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTask() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("error creating task repo: %v", err)
	}

	emptyTask := core.NewTask("", "")
	basicTask := core.NewTask("task with data", "task description")

	type args struct {
		r TaskRepo
		t *core.Task
	}
	tests := []struct {
		name    string
		args    args
		want    *core.Task
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
