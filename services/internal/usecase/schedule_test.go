package usecase_test

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	data "github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestGetSchedule(t *testing.T) {
	r := data.NewScheduleRepo()
	hourFreq, _ := schedule.NewHourFrequency([]int{0})
	hourSched := schedule.New(hourFreq, user.ID{})
	hourSchedID, _ := r.Add(hourSched)

	r2 := data.NewScheduleRepo()
	f, _ := schedule.NewHourFrequency([]int{0})
	s1 := schedule.New(f, user.ID{})
	s1.Remove()
	sID1, _ := r2.Add(s1)

	type args struct {
		r  ScheduleRepo
		id ScheduleID
	}
	tests := []struct {
		name    string
		args    args
		want    *ScheduleData
		wantErr ErrorCode
	}{
		{
			name:    "should get schedule 1",
			args:    args{r: r, id: hourSchedID},
			want:    &ScheduleData{ScheduleID: hourSchedID, Schedule: hourSched},
			wantErr: ErrNone,
		},
		{
			name:    "should return 'not found' error",
			args:    args{r: r, id: 9999},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should return 'not found' error if schedule has been removed",
			args:    args{r: r2, id: sID1},
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSchedule(tt.args.r, tt.args.id)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSchedule() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("GetSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestListSchedules(t *testing.T) {
	r1 := data.NewScheduleRepo()

	r2 := data.NewScheduleRepo()
	hourFreq, _ := schedule.NewHourFrequency([]int{0})
	hourSched := schedule.New(hourFreq, user.ID{})
	hourSchedID1, _ := r2.Add(hourSched)
	hourSchedID2, _ := r2.Add(hourSched)

	r3 := data.NewScheduleRepo()
	f, _ := schedule.NewHourFrequency([]int{0})
	s1 := schedule.New(f, user.ID{})
	s2 := schedule.New(f, user.ID{})
	s2.Remove()
	sID1, _ := r3.Add(s1)
	r3.Add(s2)

	type args struct {
		r ScheduleRepo
	}
	tests := []struct {
		name    string
		args    args
		want    map[ScheduleID]*schedule.Schedule
		wantErr ErrorCode
	}{
		{
			name:    "should return empty list",
			args:    args{r1},
			want:    map[ScheduleID]*schedule.Schedule{},
			wantErr: ErrNone,
		},
		{
			name:    "should list 2 schedules",
			args:    args{r2},
			want:    map[ScheduleID]*schedule.Schedule{hourSchedID1: hourSched, hourSchedID2: hourSched},
			wantErr: ErrNone,
		},
		{
			name:    "should list 1 schedule",
			args:    args{r3},
			want:    map[ScheduleID]*schedule.Schedule{sID1: s1},
			wantErr: ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListSchedules(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSchedules() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("ListSchedules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAddSchedule(t *testing.T) {
	r := data.NewScheduleRepo()
	hourFreq, _ := schedule.NewHourFrequency([]int{0})
	hourSchedule := schedule.New(hourFreq, user.ID{})
	c := make(chan<- bool)

	type args struct {
		r ScheduleRepo
		s *schedule.Schedule
	}
	tests := []struct {
		name    string
		args    args
		want    ScheduleID
		wantErr ErrorCode
	}{
		{
			name:    "should add an hour schedule with ID 1",
			args:    args{r: r, s: hourSchedule},
			want:    ScheduleID(1),
			wantErr: ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddSchedule(tt.args.r, tt.args.s, c)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddSchedule() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("AddSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPauseSchedule(t *testing.T) {
	r := data.NewScheduleRepo()
	f, _ := schedule.NewHourFrequency([]int{0})
	s1 := schedule.New(f, user.ID{})
	s2 := schedule.New(f, user.ID{})
	s2.Remove()
	sID1, _ := r.Add(s1)
	sID2, _ := r.Add(s2)
	c := make(chan<- bool)

	type args struct {
		r  ScheduleRepo
		id ScheduleID
	}
	tests := []struct {
		name    string
		args    args
		wantErr ErrorCode
	}{
		{
			name:    "should pause schedule",
			args:    args{r: r, id: sID1},
			wantErr: ErrNone,
		},
		{
			name:    "should return 'not found' error",
			args:    args{r: r, id: 9999},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should return 'not found' error if schedule has been removed",
			args:    args{r: r, id: sID2},
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := PauseSchedule(tt.args.r, tt.args.id, c)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("PauseSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUnpauseSchedule(t *testing.T) {
	r := data.NewScheduleRepo()
	f, _ := schedule.NewHourFrequency([]int{0})
	s1 := schedule.New(f, user.ID{})
	s2 := schedule.New(f, user.ID{})
	s1.Pause()
	s2.Remove()
	sID1, _ := r.Add(s1)
	sID2, _ := r.Add(s2)
	c := make(chan<- bool)

	type args struct {
		r  ScheduleRepo
		id ScheduleID
	}
	tests := []struct {
		name    string
		args    args
		wantErr ErrorCode
	}{
		{
			name:    "should unpause schedule",
			args:    args{r: r, id: sID1},
			wantErr: ErrNone,
		},
		{
			name:    "should return 'not found' error",
			args:    args{r: r, id: 9999},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should return 'not found' error if schedule has been removed",
			args:    args{r: r, id: sID2},
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnpauseSchedule(tt.args.r, tt.args.id, c)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UnpauseSchedule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRemoveSchedule(t *testing.T) {
	r := data.NewScheduleRepo()
	f, err := schedule.NewHourFrequency([]int{0})
	if err != nil {
		t.Fatal(err)
	}
	s := schedule.New(f, user.ID{})
	sID, err := r.Add(s)
	if err != nil {
		t.Fatal(err)
	}
	c := make(chan<- bool)

	err = RemoveSchedule(r, sID, c)
	if err != nil {
		t.Errorf("RemoveSchedule() error = %v, wantErr %v", err, nil)
	}

	retrievedS, ucErr := r.Get(sID)
	if ucErr != nil {
		t.Errorf("Get(%v) should have retrieved the removed schedule", sID)
	}
	if retrievedS.IsValid() {
		t.Errorf("RemoveSchedule() should cause IsValid() to return false for schedule: %v", retrievedS)
	}
}

func TestAddRecurringTask(t *testing.T) {
	r := data.NewScheduleRepo()
	hourFreq, _ := schedule.NewHourFrequency([]int{0})
	hourSched := schedule.New(hourFreq, user.ID{})
	hourSchedID, _ := r.Add(hourSched)
	rt1 := schedule.NewRecurringTask("task 1", "")
	rt2 := schedule.NewRecurringTask("task 2", "")

	type args struct {
		r  ScheduleRepo
		id ScheduleID
		rt schedule.RecurringTask
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr ErrorCode
	}{
		{
			name:    "should add 1st recurring task",
			args:    args{r: r, id: hourSchedID, rt: rt1},
			wantErr: ErrNone,
		},
		{
			name:    "should add 2nd recurring task",
			args:    args{r: r, id: hourSchedID, rt: rt2},
			wantErr: ErrNone,
		},
		{
			name:    "should return schedule not found error",
			args:    args{r: r, id: 9999, rt: rt1},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should return duplicate error attempting to add a duplicate recurring task",
			args:    args{r: r, id: hourSchedID, rt: rt2},
			wantErr: ErrDuplicateRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddRecurringTask(tt.args.r, tt.args.id, tt.args.rt)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("AddRecurringTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRemoveRecurringTask(t *testing.T) {
	r := data.NewScheduleRepo()
	hourFreq, _ := schedule.NewHourFrequency([]int{0})
	hourSched := schedule.New(hourFreq, user.ID{})
	rt1 := schedule.NewRecurringTask("task 1", "")
	rt2 := schedule.NewRecurringTask("task 2", "")
	rt2remove := schedule.NewRecurringTask("task 2", "")
	rt3unknown := schedule.NewRecurringTask("unknown task", "")
	hourSched.AddTask(rt1)
	hourSched.AddTask(rt2)
	hourSchedID, _ := r.Add(hourSched)

	type args struct {
		r  ScheduleRepo
		id ScheduleID
		t  schedule.RecurringTask
	}
	tests := []struct {
		name    string
		args    args
		wantErr ErrorCode
	}{
		{
			name:    "should return schedule not found error",
			args:    args{r: r, id: 9999, t: rt1},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should remove recurring task 1",
			args:    args{r: r, id: hourSchedID, t: rt1},
			wantErr: ErrNone,
		},
		{
			name:    "should remove recurring task 2",
			args:    args{r: r, id: hourSchedID, t: rt2remove},
			wantErr: ErrNone,
		},
		{
			name:    "should error attempting to remove recurring task 2 again",
			args:    args{r: r, id: hourSchedID, t: rt2},
			wantErr: ErrRecordNotFound,
		},
		{
			name:    "should error attempting to remove unknown task",
			args:    args{r: r, id: hourSchedID, t: rt3unknown},
			wantErr: ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RemoveRecurringTask(tt.args.r, tt.args.id, tt.args.t)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("RemoveRecurringTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
