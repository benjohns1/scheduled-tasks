package schedule

import (
	"reflect"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
)

func TestSchedule_Times_EmptySchedule(t *testing.T) {

	jan1st1999Midnight := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan1st2000Midnight := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan1st9999Midnight := time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC)

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
			s:       New(Frequency{}, user.ID{}),
			args:    args{jan1st9999Midnight, jan1st2000Midnight},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty schedule should return an empty slice",
			s:       New(Frequency{}, user.ID{}),
			args:    args{jan1st1999Midnight, jan1st2000Midnight},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_HourlyFrequencyOffsetInterval(t *testing.T) {

	now := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)

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
			name: "empty hour schedule should return empty slice",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 12, 0)},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "hour schedule should return 1 time",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{30})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.Add(1 * time.Hour)},
			want:    []time.Time{now.Add(30 * time.Minute)},
			wantErr: false,
		},
		{
			name: "hour schedule should return 2 times",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{30})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.Add(2 * time.Hour)},
			want:    []time.Time{now.Add(30 * time.Minute), now.Add(1*time.Hour + 30*time.Minute)},
			wantErr: false,
		},
		{
			name: "hour schedule with interval 2 should return 2 times",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{30})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				f.SetInterval(2)
				return New(f, user.ID{})
			})(),
			args:    args{now, now.Add(3 * time.Hour)},
			want:    []time.Time{now.Add(30 * time.Minute), now.Add(2*time.Hour + 30*time.Minute)},
			wantErr: false,
		},
		{
			name: "hour schedule with offset 1 should return 1 time",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{30})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				f.SetOffset(1)
				return New(f, user.ID{})
			})(),
			args:    args{now, now.Add(2 * time.Hour)},
			want:    []time.Time{now.Add(1*time.Hour + 30*time.Minute)},
			wantErr: false,
		},
		{
			name: "hour schedule with interval 2 and offset 1 should return 2 times",
			s: (func() *Schedule {
				f, err := NewHourFrequency([]int{30})
				if err != nil {
					t.Fatalf("Error creating hour frequency")
				}
				f.SetInterval(2)
				f.SetOffset(1)
				return New(f, user.ID{})
			})(),
			args:    args{now, now.Add(4 * time.Hour)},
			want:    []time.Time{now.Add(1*time.Hour + 30*time.Minute), now.Add(3*time.Hour + 30*time.Minute)},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_HourlyFrequencyEveryHour(t *testing.T) {

	jan1st2000Midnight := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	dec31st1999ElevenPM := time.Date(1999, time.December, 31, 23, 0, 0, 0, time.UTC)
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	dec31st1999ElevenPMInBeijing := time.Date(1999, time.December, 31, 23, 0, 0, 0, beijing)

	everyHourOnTheHour, _ := NewHourFrequency([]int{0})
	everyHourOnThirtyMinuteMark, _ := NewHourFrequency([]int{30})
	everyHalfHour, _ := NewHourFrequency([]int{0, 30})
	everyFifteenMinutes, _ := NewHourFrequency([]int{0, 15, 30, 45})

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
			name:    "the same start/end for schedule on the hour should return a slice with exactly 1 event on the hour",
			s:       New(everyHourOnTheHour, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{jan1st2000Midnight},
			wantErr: false,
		},
		{
			name:    "the same start/end for schedule every hour on the thirty minute mark should return an empty slice",
			s:       New(everyHourOnThirtyMinuteMark, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "should return a slice with 1 event in between boundaries",
			s:       New(everyHourOnTheHour, user.ID{}),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * -1), dec31st1999ElevenPM.Add(time.Minute)},
			want:    []time.Time{dec31st1999ElevenPM},
			wantErr: false,
		},
		{
			name:    "should return a slice with 2 events at the included boundaries of start and end",
			s:       New(everyHourOnTheHour, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour)},
			want:    []time.Time{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 3 events at the included boundaries of start and end every half hour",
			s:       New(everyHalfHour, user.ID{}),
			args:    args{jan1st2000Midnight.Add(time.Minute * 30), jan1st2000Midnight.Add(time.Minute * 90)},
			want:    []time.Time{jan1st2000Midnight.Add(time.Minute * 30), jan1st2000Midnight.Add(time.Hour), jan1st2000Midnight.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 3 events at the included boundaries of start/end before/after midnight 1999 every hour",
			s:       New(everyHourOnTheHour, user.ID{}),
			args:    args{dec31st1999ElevenPM, dec31st1999ElevenPM.Add(time.Hour * 2)},
			want:    []time.Time{dec31st1999ElevenPM, dec31st1999ElevenPM.Add(time.Hour), dec31st1999ElevenPM.Add(time.Hour * 2)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 5 events excluding boundaries before/after midnight 1999 every 15 minutes",
			s:       New(everyFifteenMinutes, user.ID{}),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * 29), dec31st1999ElevenPM.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPM.Add(time.Minute * 30), dec31st1999ElevenPM.Add(time.Minute * 45), dec31st1999ElevenPM.Add(time.Minute * 60), dec31st1999ElevenPM.Add(time.Minute * 75), dec31st1999ElevenPM.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 5 events excluding boundaries before/after midnight 1999 every 15 minutes in Beijing",
			s:       New(everyFifteenMinutes, user.ID{}),
			args:    args{dec31st1999ElevenPMInBeijing.Add(time.Minute * 29), dec31st1999ElevenPMInBeijing.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPMInBeijing.Add(time.Minute * 30), dec31st1999ElevenPMInBeijing.Add(time.Minute * 45), dec31st1999ElevenPMInBeijing.Add(time.Minute * 60), dec31st1999ElevenPMInBeijing.Add(time.Minute * 75), dec31st1999ElevenPMInBeijing.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 1 event that uses the start timezone if the end timezone differs",
			s:       New(everyHourOnTheHour, user.ID{}),
			args:    args{dec31st1999ElevenPMInBeijing.Add(time.Minute * -1), dec31st1999ElevenPM.Add(time.Second * -1 * time.Duration(secondsEastOfUTC))},
			want:    []time.Time{dec31st1999ElevenPMInBeijing},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_HourlyFrequencyEvenHour(t *testing.T) {

	jan1st2000Midnight := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	dec31st1999TenPM := time.Date(1999, time.December, 31, 22, 0, 0, 0, time.UTC)
	dec31st1999ElevenPM := time.Date(1999, time.December, 31, 23, 0, 0, 0, time.UTC)
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	dec31st1999ElevenPMInBeijing := time.Date(1999, time.December, 31, 23, 0, 0, 0, beijing)

	evenHour, _ := NewHourFrequency([]int{0})
	evenHour.SetInterval(2)
	evenHourOnThirtyMinuteMark, _ := NewHourFrequency([]int{30})
	evenHourOnThirtyMinuteMark.SetInterval(2)
	evenHourHalfHour, _ := NewHourFrequency([]int{0, 30})
	evenHourHalfHour.SetInterval(2)
	evenHourFifteenMinutes, _ := NewHourFrequency([]int{0, 15, 30, 45})
	evenHourFifteenMinutes.SetInterval(2)

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
			name:    "even hour with midnight to midnight should return a slice with 1 event",
			s:       New(evenHour, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{jan1st2000Midnight},
			wantErr: false,
		},
		{
			name:    "even hour 30 minute mark with midnight to midnight should return an empty slice",
			s:       New(evenHourOnThirtyMinuteMark, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "even hour with just before 11pm to just after 11pm should return an empty slice",
			s:       New(evenHour, user.ID{}),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * -1), dec31st1999ElevenPM.Add(time.Minute)},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "even hour with just before 10pm to just after 10pm should return an slice with 1 event",
			s:       New(evenHour, user.ID{}),
			args:    args{dec31st1999TenPM.Add(time.Minute * -1), dec31st1999TenPM.Add(time.Minute)},
			want:    []time.Time{dec31st1999TenPM},
			wantErr: false,
		},
		{
			name:    "even hour with midnight to 1am should return a slice with 1 event",
			s:       New(evenHour, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour)},
			want:    []time.Time{jan1st2000Midnight},
			wantErr: false,
		},
		{
			name:    "even hour with midnight to 2am should return a slice with 2 events",
			s:       New(evenHour, user.ID{}),
			args:    args{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour * 2)},
			want:    []time.Time{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour * 2)},
			wantErr: false,
		},
		{
			name:    "even hour every half hour with 12:30am to 2am should return a slice with 2 events",
			s:       New(evenHourHalfHour, user.ID{}),
			args:    args{jan1st2000Midnight.Add(time.Minute * 30), jan1st2000Midnight.Add(time.Hour * 2)},
			want:    []time.Time{jan1st2000Midnight.Add(time.Minute * 30), jan1st2000Midnight.Add(time.Hour * 2)},
			wantErr: false,
		},
		{
			name:    "even hour with 11pm to 1am should return a slice with 1 event",
			s:       New(evenHour, user.ID{}),
			args:    args{dec31st1999ElevenPM, dec31st1999ElevenPM.Add(time.Hour * 2)},
			want:    []time.Time{dec31st1999ElevenPM.Add(time.Hour)},
			wantErr: false,
		},
		{
			name:    "even hour every fifteen minutes with 11:29pm to 00:31am should return a slice with 3 events",
			s:       New(evenHourFifteenMinutes, user.ID{}),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * 29), dec31st1999ElevenPM.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPM.Add(time.Minute * 60), dec31st1999ElevenPM.Add(time.Minute * 75), dec31st1999ElevenPM.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "even hour every fifteen minutes with 11:29pm in Beijing to 00:31am in Beijing should return a slice with 3 events",
			s:       New(evenHourFifteenMinutes, user.ID{}),
			args:    args{dec31st1999ElevenPMInBeijing.Add(time.Minute * 29), dec31st1999ElevenPMInBeijing.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPMInBeijing.Add(time.Minute * 60), dec31st1999ElevenPMInBeijing.Add(time.Minute * 75), dec31st1999ElevenPMInBeijing.Add(time.Minute * 90)},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_DailyFrequency(t *testing.T) {

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	dec31st := time.Date(2000, time.December, 31, 0, 0, 0, 0, time.UTC)

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
			name: "empty day schedule should return empty slice",
			s: (func() *Schedule {
				f, err := NewDayFrequency([]int{}, []int{})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 12, 0)},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "day schedule should return 2 times",
			s: (func() *Schedule {
				f, err := NewDayFrequency([]int{0}, []int{1})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 0, 2)},
			want:    []time.Time{now.Add(1 * time.Hour), now.AddDate(0, 0, 1).Add(1 * time.Hour)},
			wantErr: false,
		},
		{
			name: "day schedule with interval 2 and offset 1 should return 2 times",
			s: (func() *Schedule {
				f, err := NewDayFrequency([]int{0}, []int{1})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				f.SetInterval(2)
				f.SetOffset(1)
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 0, 4)},
			want:    []time.Time{now.AddDate(0, 0, 1).Add(1 * time.Hour), now.AddDate(0, 0, 3).Add(1 * time.Hour)},
			wantErr: false,
		},
		{
			name: "day schedule with interval 3 and offset 2 should return 4 times",
			s: (func() *Schedule {
				f, err := NewDayFrequency([]int{0}, []int{6, 12})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				f.SetInterval(3)
				f.SetOffset(2)
				return New(f, user.ID{})
			})(),
			args:    args{dec31st, dec31st.AddDate(0, 0, 8)},
			want:    []time.Time{dec31st.AddDate(0, 0, 4).Add(6 * time.Hour), dec31st.AddDate(0, 0, 4).Add(12 * time.Hour), dec31st.AddDate(0, 0, 7).Add(6 * time.Hour), dec31st.AddDate(0, 0, 7).Add(12 * time.Hour)},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_WeeklyFrequency(t *testing.T) {

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

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
			name: "empty week schedule should return empty slice",
			s: (func() *Schedule {
				f, err := NewWeekFrequency([]int{}, []int{}, []time.Weekday{})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 12, 0)},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "week schedule should return 2 times",
			s: (func() *Schedule {
				f, err := NewWeekFrequency([]int{0}, []int{1}, []time.Weekday{time.Saturday})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 0, 8)},
			want:    []time.Time{now.Add(1 * time.Hour), now.AddDate(0, 0, 7).Add(1 * time.Hour)},
			wantErr: false,
		},
		{
			name: "week schedule with interval 2 and offset 1 should return 2 times",
			s: (func() *Schedule {
				f, err := NewWeekFrequency([]int{0}, []int{1}, []time.Weekday{time.Saturday})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				f.SetInterval(2)
				f.SetOffset(1)
				return New(f, user.ID{})
			})(),
			args:    args{now, time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{time.Date(2000, time.January, 15, 1, 0, 0, 0, time.UTC), time.Date(2000, time.January, 29, 1, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name: "week schedule should return 2 times",
			s: (func() *Schedule {
				f, err := NewWeekFrequency([]int{0}, []int{0}, []time.Weekday{time.Saturday, time.Sunday, time.Wednesday})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{time.Date(1999, time.December, 19, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 7, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{time.Date(1999, time.December, 19, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 22, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 25, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 26, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 29, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 2, 0, 0, 0, 0, time.UTC), time.Date(2000, time.January, 5, 0, 0, 0, 0, time.UTC)},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_Times_MonthlyFrequency(t *testing.T) {

	now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

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
			name: "empty month schedule should return empty slice",
			s: (func() *Schedule {
				f, err := NewMonthFrequency([]int{}, []int{}, []int{})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{now, now.AddDate(0, 12, 0)},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name: "month schedule should return 2 times",
			s: (func() *Schedule {
				f, err := NewMonthFrequency([]int{0}, []int{0}, []int{1})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				return New(f, user.ID{})
			})(),
			args:    args{time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name: "month schedule with interval 2 and offset 1 should return 2 times",
			s: (func() *Schedule {
				f, err := NewMonthFrequency([]int{0}, []int{0}, []int{1})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				f.SetInterval(2)
				f.SetOffset(1)
				return New(f, user.ID{})
			})(),
			args:    args{time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.May, 1, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{time.Date(2000, time.February, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.April, 1, 0, 0, 0, 0, time.UTC)},
			wantErr: false,
		},
		{
			name: "month schedule with interval 3 and offset 2 should return 4 times",
			s: (func() *Schedule {
				f, err := NewMonthFrequency([]int{0}, []int{0}, []int{1, 15, 31})
				if err != nil {
					t.Fatalf("Error creating frequency")
				}
				f.SetInterval(3)
				f.SetOffset(2)
				return New(f, user.ID{})
			})(),
			args:    args{time.Date(1999, time.November, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.March, 20, 0, 0, 0, 0, time.UTC)},
			want:    []time.Time{time.Date(1999, time.December, 1, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 15, 0, 0, 0, 0, time.UTC), time.Date(1999, time.December, 31, 0, 0, 0, 0, time.UTC), time.Date(2000, time.March, 1, 0, 0, 0, 0, time.UTC), time.Date(2000, time.March, 15, 0, 0, 0, 0, time.UTC)},
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
			if !tt.wantErr {
				nextTime, err := tt.s.NextTime(tt.args.start)
				if err != nil {
					t.Errorf("Schedule.NextTime() error = %v", err)
				}
				if len(got) > 0 && got[0] != nextTime {
					t.Errorf("Schedule.NextTime() = %v, want %v", nextTime, got[0])
				}
			}
		})
	}
}

func TestSchedule_AddTask(t *testing.T) {

	f, _ := NewHourFrequency([]int{0})
	s := New(f, user.ID{})
	rt1 := NewRecurringTask("task 1", "")
	rt2 := NewRecurringTask("task 2", "")

	type args struct {
		rt RecurringTask
	}
	tests := []struct {
		name    string
		s       *Schedule
		args    args
		wantErr bool
	}{
		{
			name:    "should add 1st recurring task",
			s:       s,
			args:    args{rt: rt1},
			wantErr: false,
		},
		{
			name:    "should add 2nd recurring task",
			s:       s,
			args:    args{rt: rt2},
			wantErr: false,
		},
		{
			name:    "should fail trying to add identical task",
			s:       s,
			args:    args{rt: rt2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.s.AddTask(tt.args.rt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Schedule.AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSchedule_RemoveTask(t *testing.T) {
	f, _ := NewHourFrequency([]int{0})
	s := New(f, user.ID{})
	rt1 := NewRecurringTask("task 1", "")
	rt2 := NewRecurringTask("task 2", "")
	rt2remove := NewRecurringTask("task 2", "")
	rt3unknown := NewRecurringTask("unknown task", "")
	s.AddTask(rt1)
	s.AddTask(rt2)

	type args struct {
		rt RecurringTask
	}
	tests := []struct {
		name    string
		s       *Schedule
		args    args
		wantErr bool
	}{
		{
			name:    "should remove recurring task 1",
			s:       s,
			args:    args{rt: rt1},
			wantErr: false,
		},
		{
			name:    "should remove recurring task 2",
			s:       s,
			args:    args{rt: rt2remove},
			wantErr: false,
		},
		{
			name:    "should error attempting to remove recurring task 2 again",
			s:       s,
			args:    args{rt: rt2},
			wantErr: true,
		},
		{
			name:    "should error attempting to remove unknown task",
			s:       s,
			args:    args{rt: rt3unknown},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.RemoveTask(tt.args.rt); (err != nil) != tt.wantErr {
				t.Errorf("Schedule.RemoveTask() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchedule_TaskListUpdates(t *testing.T) {
	f, _ := NewHourFrequency([]int{0})
	s := New(f, user.ID{})
	rt1 := NewRecurringTask("task 1", "")
	rt2 := NewRecurringTask("task 2", "")
	rt2remove := NewRecurringTask("task 2", "")
	rt3unknown := NewRecurringTask("unknown task", "")
	s.AddTask(rt1)
	s.AddTask(rt2)
	s.RemoveTask(rt2remove)
	s.RemoveTask(rt3unknown)

	t.Run("should list first task", func(t *testing.T) {
		want := []RecurringTask{rt1}
		got := s.Tasks()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Schedule.List() = %v, want %v", got, want)
		}
	})
}

func TestSchedule_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    *Schedule
		want bool
	}{
		{
			name: "schedule with a zero removed time should be valid",
			s:    &Schedule{removedTime: time.Time{}},
			want: true,
		},
		{
			name: "schedule with a removed time should be invalid",
			s:    &Schedule{removedTime: time.Now()},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsValid(); got != tt.want {
				t.Errorf("Schedule.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
