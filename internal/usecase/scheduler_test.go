package usecase_test

import (
	"reflect"
	"testing"
	"time"

	data "github.com/benjohns1/scheduled-tasks/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/internal/usecase"
)

func TestCheckSchedules(t *testing.T) {
	type args struct {
		taskRepo     TaskRepo
		scheduleRepo ScheduleRepo
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "should return zero time for next run from empty schedules",
			args: args{
				taskRepo:     data.NewTaskRepo(),
				scheduleRepo: data.NewScheduleRepo(),
			},
			want:    time.Time{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckSchedules(tt.args.taskRepo, tt.args.scheduleRepo)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSchedules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckSchedules() = %v, want %v", got, tt.want)
			}
		})
	}
}
