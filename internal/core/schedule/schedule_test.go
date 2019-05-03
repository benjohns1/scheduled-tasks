package schedule

import (
	"reflect"
	"testing"
	"time"
)

func TestSchedule_Times_HourlyFrequency(t *testing.T) {

	jan1st1999Midnight := time.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan1st2000Midnight := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	jan1st2000MidnightThirty := time.Date(2000, time.January, 1, 0, 30, 0, 0, time.UTC)
	jan1st9999Midnight := time.Date(9999, time.January, 1, 0, 0, 0, 0, time.UTC)

	dec31st1999ElevenPM := time.Date(1999, time.December, 31, 23, 0, 0, 0, time.UTC)
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
	dec31st1999ElevenPMInBeijing := time.Date(1999, time.December, 31, 23, 0, 0, 0, beijing)

	everyHourOnTheHour, _ := NewHourFrequency(1, []int{0})
	everyHourOnThirtyMinuteMark, _ := NewHourFrequency(1, []int{30})
	everyHalfHour, _ := NewHourFrequency(1, []int{0, 30})
	everyFifteenMinutes, _ := NewHourFrequency(1, []int{0, 15, 30, 45})

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
			args:    args{jan1st9999Midnight, jan1st2000Midnight},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty schedule should return an empty slice",
			s:       New(),
			args:    args{jan1st1999Midnight, jan1st2000Midnight},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "the same start/end for schedule on the hour should return a slice with exactly 1 event on the hour",
			s:       New().WithFrequency(everyHourOnTheHour),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{jan1st2000Midnight},
			wantErr: false,
		},
		{
			name:    "the same start/end for schedule every hour on the thirty minute mark should return an empty slice",
			s:       New().WithFrequency(everyHourOnThirtyMinuteMark),
			args:    args{jan1st2000Midnight, jan1st2000Midnight},
			want:    []time.Time{},
			wantErr: false,
		},
		{
			name:    "should return a slice with 1 event in between boundaries",
			s:       New().WithFrequency(everyHourOnTheHour),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * -1), dec31st1999ElevenPM.Add(time.Minute)},
			want:    []time.Time{dec31st1999ElevenPM},
			wantErr: false,
		},
		{
			name:    "should return a slice with 2 events at the included boundaries of start and end",
			s:       New().WithFrequency(everyHourOnTheHour),
			args:    args{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour)},
			want:    []time.Time{jan1st2000Midnight, jan1st2000Midnight.Add(time.Hour)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 3 events at the included boundaries of start and end every half hour",
			s:       New().WithFrequency(everyHalfHour),
			args:    args{jan1st2000MidnightThirty, jan1st2000MidnightThirty.Add(time.Hour)},
			want:    []time.Time{jan1st2000MidnightThirty, jan1st2000MidnightThirty.Add(time.Minute * 30), jan1st2000MidnightThirty.Add(time.Hour)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 3 events at the included boundaries of start/end before/after midnight 1999 every hour",
			s:       New().WithFrequency(everyHourOnTheHour),
			args:    args{dec31st1999ElevenPM, dec31st1999ElevenPM.Add(time.Hour * 2)},
			want:    []time.Time{dec31st1999ElevenPM, dec31st1999ElevenPM.Add(time.Hour), dec31st1999ElevenPM.Add(time.Hour * 2)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 5 events excluding boundaries before/after midnight 1999 every 15 minutes",
			s:       New().WithFrequency(everyFifteenMinutes),
			args:    args{dec31st1999ElevenPM.Add(time.Minute * 29), dec31st1999ElevenPM.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPM.Add(time.Minute * 30), dec31st1999ElevenPM.Add(time.Minute * 45), dec31st1999ElevenPM.Add(time.Minute * 60), dec31st1999ElevenPM.Add(time.Minute * 75), dec31st1999ElevenPM.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 5 events excluding boundaries before/after midnight 1999 every 15 minutes in Beijing",
			s:       New().WithFrequency(everyFifteenMinutes),
			args:    args{dec31st1999ElevenPMInBeijing.Add(time.Minute * 29), dec31st1999ElevenPMInBeijing.Add(time.Minute * 91)},
			want:    []time.Time{dec31st1999ElevenPMInBeijing.Add(time.Minute * 30), dec31st1999ElevenPMInBeijing.Add(time.Minute * 45), dec31st1999ElevenPMInBeijing.Add(time.Minute * 60), dec31st1999ElevenPMInBeijing.Add(time.Minute * 75), dec31st1999ElevenPMInBeijing.Add(time.Minute * 90)},
			wantErr: false,
		},
		{
			name:    "should return a slice with 1 event that uses the start timezone if the end timezone differs",
			s:       New().WithFrequency(everyHourOnTheHour),
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
		})
	}
}
