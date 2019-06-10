package task

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
)

func TestNew(t *testing.T) {
	testNow := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(testNow)
	clock.Set(clockMock)
	defer clock.Set(prevClock)

	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *Task
	}{
		{
			name: "create new task",
			args: args{name: "task name", description: "task description"},
			want: &Task{name: "task name", description: "task description", completedTime: time.Time{}, clearedTime: time.Time{}, createdTime: testNow},
		},
		{
			name: "create new empty task",
			args: args{name: "", description: ""},
			want: &Task{name: "", description: "", completedTime: time.Time{}, clearedTime: time.Time{}, createdTime: testNow},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name, tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRaw(t *testing.T) {
	type args struct {
		name        string
		description string
		complete    time.Time
		cleared     time.Time
		created     time.Time
	}
	now := time.Now()
	tests := []struct {
		name string
		args args
		want *Task
	}{
		{
			name: "create new task with all params",
			args: args{name: "task name", description: "task description", complete: now, cleared: now, created: now},
			want: &Task{name: "task name", description: "task description", completedTime: now, clearedTime: now, createdTime: now},
		},
		{
			name: "create new empty task with all params",
			args: args{name: "", description: "", complete: time.Time{}, cleared: time.Time{}, created: time.Time{}},
			want: &Task{name: "", description: "", completedTime: time.Time{}, clearedTime: time.Time{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRaw(tt.args.name, tt.args.description, tt.args.complete, tt.args.cleared, tt.args.created); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsValid(t *testing.T) {
	tests := []struct {
		name string
		t    *Task
		want bool
	}{
		{
			name: "task with a zero cleared time should be valid",
			t:    &Task{clearedTime: time.Time{}},
			want: true,
		},
		{
			name: "task with a cleared time should be invalid",
			t:    &Task{clearedTime: time.Now()},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.want {
				t.Errorf("Task.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_CompleteNow(t *testing.T) {
	tests := []struct {
		name    string
		t       *Task
		want    bool
		wantErr bool
	}{
		{
			name:    "incomplete task should be completed",
			t:       &Task{},
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid task should return error",
			t:       &Task{clearedTime: time.Now()},
			want:    false,
			wantErr: true,
		},
		{
			name:    "already completed task should not be completed again ",
			t:       &Task{completedTime: time.Now()},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.CompleteNow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.CompleteNow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Task.CompleteNow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_ClearCompleted(t *testing.T) {
	tests := []struct {
		name    string
		t       *Task
		wantErr bool
	}{
		{
			name:    "incomplete task should not be cleared",
			t:       &Task{},
			wantErr: true,
		},
		{
			name:    "completed task should be cleared",
			t:       &Task{completedTime: time.Now()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.ClearCompleted(); (err != nil) != tt.wantErr {
				t.Errorf("Task.ClearCompleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
