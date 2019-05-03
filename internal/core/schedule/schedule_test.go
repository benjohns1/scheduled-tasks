package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestSchedule_Times(t *testing.T) {
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name    string
		s       *Schedule
		args    args
		want    []time.Time
		wantErr bool
	}{
		{
			name:    "should return an error if end time is before start time",
			s:       New(),
			args:    args{time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty schedule should return an empty array",
			s:       New(),
			args:    args{time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Times(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schedule.Times() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Schedule.Times() = %v, want %v", got, tt.want)
			}
		})
	}
}
