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

	userRepo := data.NewUserRepo()
	u1 := user.New("test user for GetTask")
	uid1 := u1.ID()
	userRepo.AddExternal(u1, "p1", "e1")
	u2 := user.New("test user 2 for GetTask")
	uid2 := u2.ID()
	userRepo.AddExternal(u2, "p1", "e2")

	taskRepo := data.NewTaskRepo()
	task1 := task.New("task1", "", uid1)
	taskID, _ := taskRepo.Add(task1)
	task2 := task.NewRaw("task2", "", now, time.Time{}, now, uid1)
	completedTaskID, _ := taskRepo.Add(task2)
	clearedTaskID, _ := taskRepo.Add(task.NewRaw("task3", "", now, now, now, uid1))

	type args struct {
		r   TaskRepo
		id  TaskID
		uid user.ID
	}
	tests := []struct {
		name    string
		args    args
		want    *TaskData
		wantErr ErrorCode
	}{
		{
			name:    "task should be retrieved",
			args:    args{taskRepo, taskID, uid1},
			want:    &TaskData{TaskID: taskID, Task: task1},
			wantErr: ErrNone,
		},
		{
			name:    "completed task should be retrieved",
			args:    args{taskRepo, completedTaskID, uid1},
			want:    &TaskData{TaskID: completedTaskID, Task: task2},
			wantErr: ErrNone,
		},
		{
			name:    "cleared task should return an ErrRecordNotFound",
			args:    args{taskRepo, clearedTaskID, uid1},
			want:    nil,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "non-existent task ID should return an ErrRecordNotFound",
			args:    args{taskRepo, 99999, uid1},
			want:    nil,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "task ID for different user should return an ErrRecordNotFound",
			args:    args{taskRepo, taskID, uid2},
			want:    nil,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTask(tt.args.r, tt.args.id, tt.args.uid)
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
	r := data.NewTaskRepo()
	u1 := user.New("new user 1 for CompleteTask")
	uid1 := u1.ID()
	taskID, _ := r.Add(task.New("task1", "", uid1))
	completedTaskID, _ := r.Add(task.NewRaw("task2", "", now, time.Time{}, now, uid1))
	clearedTaskID, _ := r.Add(task.NewRaw("task3", "", now, now, now, uid1))

	u2 := user.New("new user 2 for CompleteTask")
	uid2 := u2.ID()
	u2t1, _ := r.Add(task.New("u2t1", "", uid2))

	type args struct {
		r   TaskRepo
		id  TaskID
		uid user.ID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr ErrorCode
	}{
		{
			name:    "task should be completed",
			args:    args{r, taskID, uid1},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "completed task should not be completed again",
			args:    args{r, completedTaskID, uid1},
			want:    false,
			wantErr: ErrNone,
		},
		{
			name:    "completing a cleared task should return an ErrRecordNotFound",
			args:    args{r, clearedTaskID, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "completing a non-existent task ID should return an ErrRecordNotFound",
			args:    args{r, 99999, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "completing a task created by another user should return an ErrRecordNotFound",
			args:    args{r, u2t1, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompleteTask(tt.args.r, tt.args.id, tt.args.uid)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("CompleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("CompleteTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClearTask(t *testing.T) {
	now := clock.Now()

	u1 := user.New("new user 1 for CompleteTask")
	uid1 := u1.ID()
	r := data.NewTaskRepo()
	taskID, _ := r.Add(task.New("task1", "", uid1))
	completedTaskID, _ := r.Add(task.NewRaw("task2", "", now, time.Time{}, now, uid1))
	clearedTaskID, _ := r.Add(task.NewRaw("task3", "", now, now, now, uid1))

	u2 := user.New("new user 2 for CompleteTask")
	uid2 := u2.ID()
	u2t1, _ := r.Add(task.New("u2t1", "", uid2))

	type args struct {
		r   TaskRepo
		id  TaskID
		uid user.ID
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr ErrorCode
	}{
		{
			name:    "task should be cleared",
			args:    args{r, taskID, uid1},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "completed task should be cleared",
			args:    args{r, completedTaskID, uid1},
			want:    true,
			wantErr: ErrNone,
		},
		{
			name:    "previously cleared task should return ErrRecordNotFound",
			args:    args{r, clearedTaskID, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "clearing a non-existent task ID should return ErrRecordNotFound",
			args:    args{r, 9999, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "clearing a task created by another user should return an ErrRecordNotFound",
			args:    args{r, u2t1, uid1},
			want:    false,
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ClearTask(tt.args.r, tt.args.id, tt.args.uid)
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

	u1 := user.New("new user 1 for ClearCompletedTasks")
	uid1 := u1.ID()

	singleCompletedTaskRepo := data.NewTaskRepo()
	singleCompletedTaskRepo.Add(task.New("task1", "", uid1))
	singleCompletedTaskRepo.Add(task.NewRaw("task2", "", now, time.Time{}, now, uid1))
	singleCompletedTaskRepo.Add(task.NewRaw("task3", "", now, now, now, uid1))

	thousandCompletedTasksRepo := data.NewTaskRepo()
	for i := 0; i < 1000; i++ {
		_, err := thousandCompletedTasksRepo.Add(task.NewRaw("", "", now, time.Time{}, now, uid1))
		if err != nil {
			t.Errorf("error setting up task repo: %v", err)
		}
	}

	mixedRepo := data.NewTaskRepo()
	u2 := user.New("new user 2 for ClearCompletedTasks")
	uid2 := u2.ID()
	u3 := user.New("new user 3 for ClearCompletedTasks")
	uid3 := u3.ID()
	u1t1 := task.New("mixed u1t1", "", uid1)
	u1t1.CompleteNow()
	mixedRepo.Add(u1t1)
	u1t2 := task.New("mixed u1t2", "", uid1)
	u1t2.CompleteNow()
	mixedRepo.Add(u1t2)
	u2t1 := task.New("mixed u2t1", "", uid2)
	u2t1.CompleteNow()
	mixedRepo.Add(u2t1)

	type args struct {
		r   TaskRepo
		uid user.ID
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name:      "clearing empty repo should return 0 count",
			args:      args{emptyRepo, uid1},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "clearing repo should return 1 count",
			args:      args{singleCompletedTaskRepo, uid1},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "clearing repo should return 1000 count",
			args:      args{thousandCompletedTasksRepo, uid1},
			wantCount: 1000,
			wantErr:   false,
		},
		{
			name:      "clearing repo should not clear other user's tasks",
			args:      args{mixedRepo, uid3},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "clearing repo should only clear user's own tasks",
			args:      args{mixedRepo, uid1},
			wantCount: 2,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := ClearCompletedTasks(tt.args.r, tt.args.uid)
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
	userRepo := data.NewUserRepo()
	u1 := user.New("test user ListTasks")
	uid1 := u1.ID()
	userRepo.AddExternal(u1, "p1", "e1")
	taskRepo := data.NewTaskRepo()
	task1 := task.New("task1", "", uid1)
	id1, _ := taskRepo.Add(task1)
	task2 := task.NewRaw("task2", "", now, time.Time{}, now, uid1)
	id2, _ := taskRepo.Add(task2)
	taskRepo.Add(task.NewRaw("task3", "", now, now, now, uid1))

	type args struct {
		r   TaskRepo
		uid user.ID
	}
	tests := []struct {
		name    string
		args    args
		want    map[TaskID]*task.Task
		wantErr ErrorCode
	}{
		{
			name: "task list should return map of 2 uncleared tasks",
			args: args{taskRepo, uid1},
			want: map[TaskID]*task.Task{
				id1: task1,
				id2: task2,
			},
			wantErr: ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTasks(tt.args.r, tt.args.uid)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ListTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}
