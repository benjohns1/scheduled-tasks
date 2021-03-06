package transient

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestNewScheduleRepo(t *testing.T) {
	tests := []struct {
		name string
		want *ScheduleRepo
	}{
		{
			name: "should return new empty repo",
			want: &ScheduleRepo{lastID: 0, schedules: make(map[usecase.ScheduleID]*schedule.Schedule)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScheduleRepo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScheduleRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScheduleRepo_Get(t *testing.T) {
	r := NewScheduleRepo()
	emptyHourlyFreq, _ := schedule.NewHourFrequency([]int{0})
	emptyHourlySched := schedule.New(emptyHourlyFreq, user.ID{})
	emptyID, _ := r.Add(emptyHourlySched)

	type args struct {
		id usecase.ScheduleID
	}
	tests := []struct {
		name    string
		r       *ScheduleRepo
		args    args
		want    *schedule.Schedule
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 1 empty hourly schedule",
			r:       r,
			args:    args{id: emptyID},
			want:    emptyHourlySched,
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Get(tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScheduleRepo.Get() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScheduleRepo_GetAll(t *testing.T) {
	r := NewScheduleRepo()
	emptyHourlyFreq1, _ := schedule.NewHourFrequency([]int{0})
	emptyHourlySched1 := schedule.New(emptyHourlyFreq1, user.ID{})
	emptyHourlyFreq2, _ := schedule.NewHourFrequency([]int{0})
	emptyHourlySched2 := schedule.New(emptyHourlyFreq2, user.ID{})
	id1, _ := r.Add(emptyHourlySched1)
	id2, _ := r.Add(emptyHourlySched2)

	tests := []struct {
		name    string
		r       *ScheduleRepo
		want    map[usecase.ScheduleID]*schedule.Schedule
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 2 empty hourly schedules",
			r:       r,
			want:    map[usecase.ScheduleID]*schedule.Schedule{id1: emptyHourlySched1, id2: emptyHourlySched2},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAll()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScheduleRepo.GetAll() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScheduleRepo_Add(t *testing.T) {
	r := NewScheduleRepo()
	emptyHourlyFreq1, _ := schedule.NewHourFrequency([]int{0})
	emptyHourlySched1 := schedule.New(emptyHourlyFreq1, user.ID{})

	type args struct {
		s *schedule.Schedule
	}
	tests := []struct {
		name    string
		r       *ScheduleRepo
		args    args
		want    usecase.ScheduleID
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should add 1 empty hourly schedule",
			r:       r,
			args:    args{s: emptyHourlySched1},
			want:    1,
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Add(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScheduleRepo.Add() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScheduleRepo_Update(t *testing.T) {
	r := NewScheduleRepo()
	hourlyFreq1, _ := schedule.NewHourFrequency([]int{0})
	hourlySched1 := schedule.New(hourlyFreq1, user.ID{})
	id1, _ := r.Add(hourlySched1)
	hourlySched1.Pause()

	type args struct {
		id usecase.ScheduleID
		s  *schedule.Schedule
	}
	tests := []struct {
		name    string
		r       *ScheduleRepo
		args    args
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should successfully update schedule",
			r:       r,
			args:    args{id: id1, s: hourlySched1},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Update(tt.args.id, tt.args.s)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
