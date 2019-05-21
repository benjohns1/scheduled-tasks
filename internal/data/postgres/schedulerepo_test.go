// +build integration

package postgres_test

import (
	"reflect"
	"testing"

	. "github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	. "github.com/benjohns1/scheduled-tasks/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

func TestNewScheduleRepo(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	type args struct {
		conn DBConn
	}
	tests := []struct {
		name          string
		args          args
		wantSchedules map[usecase.ScheduleID]*schedule.Schedule
		wantErr       bool
	}{
		{
			name:          "should return new empty repo",
			args:          args{conn},
			wantSchedules: map[usecase.ScheduleID]*schedule.Schedule{},
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepo, err := NewScheduleRepo(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewScheduleRepo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotSchedules, err := gotRepo.GetAll()
			if err != nil {
				t.Errorf("NewScheduleRepo() error retrieving schedules: %v", err)
			}
			if !reflect.DeepEqual(gotSchedules, tt.wantSchedules) {
				t.Errorf("NewScheduleRepo() schedules = %v, want %v", gotSchedules, tt.wantSchedules)
			}
		})
	}
}

func TestScheduleRepo_Get(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	emptyHourlyFreq, err := schedule.NewHourFrequency([]int{0})
	if err != nil {
		t.Fatal(err)
	}
	emptyHourlySched := schedule.New(emptyHourlyFreq)
	emptyID, err := r.Add(emptyHourlySched)
	if err != nil {
		t.Fatal(err)
	}

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
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	emptyFreq1 := schedule.Frequency{}
	if err != nil {
		t.Fatal(err)
	}
	emptySched1 := schedule.New(emptyFreq1)
	emptyFreq2 := schedule.Frequency{}
	if err != nil {
		t.Fatal(err)
	}
	emptySched2 := schedule.New(emptyFreq2)
	id1, err := r.Add(emptySched1)
	if err != nil {
		t.Fatal(err)
	}
	id2, err := r.Add(emptySched2)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		r       *ScheduleRepo
		wantMap map[usecase.ScheduleID]*schedule.Schedule
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 2 empty hourly schedules",
			r:       r,
			wantMap: map[usecase.ScheduleID]*schedule.Schedule{id1: emptySched1, id2: emptySched2},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAll()
			if len(got) != len(tt.wantMap) {
				t.Errorf("ScheduleRepo.GetAll() got = %v, want %v", got, tt.wantMap)
			}
			for id, schedule := range got {
				if !reflect.DeepEqual(schedule, tt.wantMap[id]) {
					t.Errorf("ScheduleRepo.GetAll() schedule[%v] got = %v, want %v", id, schedule, tt.wantMap[id])
				}
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}


func TestScheduleRepo_GetAllUnpaused(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	f, err := schedule.NewHourFrequency([]int{})
	if err != nil {
		t.Fatal(err)
	}
	s1 := schedule.New(f)
	s1.Pause()

	s2 := schedule.New(f)

	_, err = r.Add(s1)
	if err != nil {
		t.Fatal(err)
	}

	id2, err := r.Add(s2)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		r       *ScheduleRepo
		wantMap map[usecase.ScheduleID]*schedule.Schedule
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get 1 empty hourly schedule",
			r:       r,
			wantMap: map[usecase.ScheduleID]*schedule.Schedule{id2: s2},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAllUnpaused()
			if len(got) != len(tt.wantMap) {
				t.Errorf("ScheduleRepo.GetAllUnpaused() got = %v, want %v", got, tt.wantMap)
			}
			for id, schedule := range got {
				if !reflect.DeepEqual(schedule, tt.wantMap[id]) {
					t.Errorf("ScheduleRepo.GetAllUnpaused() schedule[%v] got = %v, want %v", id, schedule, tt.wantMap[id])
				}
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.GetAllUnpaused() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestScheduleRepo_Add(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	emptyHourlyFreq1, err := schedule.NewHourFrequency([]int{0})
	if err != nil {
		t.Fatal(err)
	}
	emptyHourlySched1 := schedule.New(emptyHourlyFreq1)

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
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	hourlyFreq1, _ := schedule.NewHourFrequency([]int{0})
	hourlySched1 := schedule.New(hourlyFreq1)
	id1, err := r.Add(hourlySched1)
	if err != nil {
		t.Fatal(err)
	}
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


func TestScheduleRepo_AnyTasksModified(t *testing.T) {
	type args struct {
		as map[int64]schedule.RecurringTask
		bs []schedule.RecurringTask
	}
	tests := []struct {
		name    string
		args    args
		want 	bool
	}{
		{
			name:    "two empty lists should return false",
			args:    args{
				as: map[int64]schedule.RecurringTask{},
				bs: []schedule.RecurringTask{},
			},
			want:    false,
		},
		{
			name:    "two equal lists with 1 task should return false",
			args:    args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc")},
			},
			want:    false,
		},
		{
			name:    "two equal lists with 2 tasks should return false",
			args:    args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc"),999: schedule.NewRecurringTask("999", "999 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"),schedule.NewRecurringTask("999", "999 desc")},
			},
			want:    false,
		},
		{
			name:    "empty map and slice with 1 task should return true",
			args:    args{
				as: map[int64]schedule.RecurringTask{},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc")},
			},
			want:    true,
		},
		{
			name:    "map with 1 task and empty slice should return true",
			args:    args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{},
			},
			want:    true,
		},
		{
			name:    "two lists with different number of tasks should return true",
			args:    args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"),schedule.NewRecurringTask("999", "999 desc")},
			},
			want:    true,
		},
		{
			name:    "two lists with differring tasks should return true",
			args:    args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc"), 2: schedule.NewRecurringTask("2", "2 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"),schedule.NewRecurringTask("999", "999 desc")},
			},
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AnyTasksModified(tt.args.as, tt.args.bs)
			if got != tt.want {
				t.Errorf("schedule.AnyTasksModified() = %v, want %v", got, tt.want)
			}
		})
	}
}
