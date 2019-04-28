package usecase_test

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	data "github.com/benjohns1/scheduled-tasks/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

func TestAddTask(t *testing.T) {
	taskRepo, err := data.NewTaskRepo()
	if err != nil {
		t.Errorf("Error getting test task repo: %v", err)
	}

	type args struct {
		r    *data.TaskRepo
		name string
		desc string
	}
	tests := []struct {
		name    string
		args    args
		want    *usecase.TaskData
		wantErr bool
	}{
		{
			name:    "add empty task",
			args:    args{r: taskRepo, name: "", desc: ""},
			want:    &usecase.TaskData{Task: core.NewTask("", "")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usecase.AddTask(tt.args.r, tt.args.name, tt.args.desc)
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
