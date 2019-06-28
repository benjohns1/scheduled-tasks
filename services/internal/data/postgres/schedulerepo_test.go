// +build integration

package postgres_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
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

func addHourSchedule(t *testing.T, r *ScheduleRepo, atMinutes []int) (f schedule.Frequency, s *schedule.Schedule, id usecase.ScheduleID) {

	f, err := schedule.NewHourFrequency(atMinutes)
	if err != nil {
		t.Fatal(err)
	}
	s, id = addSchedule(t, r, f)
	return f, s, id
}

func addDaySchedule(t *testing.T, r *ScheduleRepo, atMinutes []int, atHours []int) (f schedule.Frequency, s *schedule.Schedule, id usecase.ScheduleID) {

	f, err := schedule.NewDayFrequency(atMinutes, atHours)
	if err != nil {
		t.Fatal(err)
	}
	s, id = addSchedule(t, r, f)
	return f, s, id
}

func addWeekSchedule(t *testing.T, r *ScheduleRepo, atMinutes []int, atHours []int, onDaysOfWeek []time.Weekday) (f schedule.Frequency, s *schedule.Schedule, id usecase.ScheduleID) {

	f, err := schedule.NewWeekFrequency(atMinutes, atHours, onDaysOfWeek)
	if err != nil {
		t.Fatal(err)
	}
	s, id = addSchedule(t, r, f)
	return f, s, id
}

func addSchedule(t *testing.T, r *ScheduleRepo, f schedule.Frequency) (s *schedule.Schedule, id usecase.ScheduleID) {
	s = schedule.New(f)
	id, err := r.Add(s)
	if err != nil {
		t.Fatal(err)
	}
	return s, id
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

	_, hs, hsID := addHourSchedule(t, r, []int{0})
	_, ds, dsID := addDaySchedule(t, r, []int{0}, []int{0})
	_, ws, wsID := addWeekSchedule(t, r, []int{0}, []int{0}, []time.Weekday{time.Sunday})

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
			name:    "should get hour schedule",
			r:       r,
			args:    args{id: hsID},
			want:    hs,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should get day schedule",
			r:       r,
			args:    args{id: dsID},
			want:    ds,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should get week schedule",
			r:       r,
			args:    args{id: wsID},
			want:    ws,
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

	f1 := schedule.Frequency{}
	if err != nil {
		t.Fatal(err)
	}
	s1, id1 := addSchedule(t, r, f1)

	f2 := schedule.Frequency{}
	if err != nil {
		t.Fatal(err)
	}
	s2, id2 := addSchedule(t, r, f2)

	_, hs, hsID := addHourSchedule(t, r, []int{0})
	_, ds, dsID := addDaySchedule(t, r, []int{0}, []int{0})
	_, ws, wsID := addWeekSchedule(t, r, []int{0}, []int{0}, []time.Weekday{time.Sunday})

	tests := []struct {
		name    string
		r       *ScheduleRepo
		wantMap map[usecase.ScheduleID]*schedule.Schedule
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get all schedules",
			r:       r,
			wantMap: map[usecase.ScheduleID]*schedule.Schedule{id1: s1, id2: s2, hsID: hs, dsID: ds, wsID: ws},
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

func TestScheduleRepo_GetAllScheduled(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	f, err := schedule.NewHourFrequency([]int{})
	if err != nil {
		t.Fatal(err)
	}
	r, err := NewScheduleRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	sPause := schedule.New(f)
	sPause.Pause()
	sRemove := schedule.New(f)
	sRemove.Remove()
	sValid := schedule.New(f)
	_, err = r.Add(sPause)
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Add(sRemove)
	if err != nil {
		t.Fatal(err)
	}
	validID, err := r.Add(sValid)
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
			name:    "should not return paused or removed schedules",
			r:       r,
			wantMap: map[usecase.ScheduleID]*schedule.Schedule{validID: sValid},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAllScheduled()
			if len(got) != len(tt.wantMap) {
				t.Errorf("ScheduleRepo.GetAllScheduled() got = %v, want %v", got, tt.wantMap)
			}
			for id, schedule := range got {
				if !reflect.DeepEqual(schedule, tt.wantMap[id]) {
					t.Errorf("ScheduleRepo.GetAllScheduled() schedule[%v] got = %v, want %v", id, schedule, tt.wantMap[id])
				}
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ScheduleRepo.GetAllScheduled() error = %v, wantErr %v", err, tt.wantErr)
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

	hf, err := schedule.NewHourFrequency([]int{0})
	if err != nil {
		t.Fatal(err)
	}
	df, err := schedule.NewDayFrequency([]int{0}, []int{0})
	if err != nil {
		t.Fatal(err)
	}
	wf, err := schedule.NewWeekFrequency([]int{0}, []int{0}, []time.Weekday{time.Sunday})
	if err != nil {
		t.Fatal(err)
	}

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
			name:    "should add hour schedule",
			r:       r,
			args:    args{s: schedule.New(hf)},
			want:    1,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should add day schedule",
			r:       r,
			args:    args{s: schedule.New(df)},
			want:    2,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should add week schedule",
			r:       r,
			args:    args{s: schedule.New(wf)},
			want:    3,
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

	hf, _ := schedule.NewHourFrequency([]int{0})
	hs := schedule.New(hf)
	hsID, err := r.Add(hs)
	if err != nil {
		t.Fatal(err)
	}
	hs.Pause()

	df, _ := schedule.NewDayFrequency([]int{0}, []int{0})
	ds := schedule.New(df)
	dsID, err := r.Add(ds)
	if err != nil {
		t.Fatal(err)
	}
	ds.Pause()

	wf, _ := schedule.NewWeekFrequency([]int{0}, []int{0}, []time.Weekday{time.Sunday})
	ws := schedule.New(wf)
	wsID, err := r.Add(ws)
	if err != nil {
		t.Fatal(err)
	}
	ws.Pause()

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
			name:    "should successfully update hour schedule",
			r:       r,
			args:    args{id: hsID, s: hs},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should successfully update day schedule",
			r:       r,
			args:    args{id: dsID, s: ds},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should successfully update week schedule",
			r:       r,
			args:    args{id: wsID, s: ws},
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
		name string
		args args
		want bool
	}{
		{
			name: "two empty lists should return false",
			args: args{
				as: map[int64]schedule.RecurringTask{},
				bs: []schedule.RecurringTask{},
			},
			want: false,
		},
		{
			name: "two equal lists with 1 task should return false",
			args: args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc")},
			},
			want: false,
		},
		{
			name: "two equal lists with 2 tasks should return false",
			args: args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc"), 999: schedule.NewRecurringTask("999", "999 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"), schedule.NewRecurringTask("999", "999 desc")},
			},
			want: false,
		},
		{
			name: "empty map and slice with 1 task should return true",
			args: args{
				as: map[int64]schedule.RecurringTask{},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc")},
			},
			want: true,
		},
		{
			name: "map with 1 task and empty slice should return true",
			args: args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{},
			},
			want: true,
		},
		{
			name: "two lists with different number of tasks should return true",
			args: args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"), schedule.NewRecurringTask("999", "999 desc")},
			},
			want: true,
		},
		{
			name: "two lists with differring tasks should return true",
			args: args{
				as: map[int64]schedule.RecurringTask{1: schedule.NewRecurringTask("1", "1 desc"), 2: schedule.NewRecurringTask("2", "2 desc")},
				bs: []schedule.RecurringTask{schedule.NewRecurringTask("1", "1 desc"), schedule.NewRecurringTask("999", "999 desc")},
			},
			want: true,
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
