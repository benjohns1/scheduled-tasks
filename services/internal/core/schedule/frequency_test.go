package schedule

import (
	"reflect"
	"testing"
	"time"
)

func emptyFrequencies(t *testing.T) []Frequency {
	hourF, err := NewHourFrequency([]int{0})
	if err != nil {
		t.Fatalf("NewHourFrequency() returned unexpected error = %v", err)
	}
	dayF, err := NewDayFrequency([]int{0}, []int{0})
	if err != nil {
		t.Fatalf("NewDayFrequency() returned unexpected error = %v", err)
	}
	weekF, err := NewWeekFrequency([]int{0}, []int{0}, []time.Weekday{time.Sunday})
	if err != nil {
		t.Fatalf("NewWeekFrequency() returned unexpected error = %v", err)
	}
	monthF, err := NewMonthFrequency([]int{0}, []int{0}, []int{1})
	if err != nil {
		t.Fatalf("NewMonthFrequency() returned unexpected error = %v", err)
	}
	return []Frequency{
		hourF,
		dayF,
		weekF,
		monthF,
	}
}

func TestDefaultOffsetInterval(t *testing.T) {
	frequencies := emptyFrequencies(t)

	defaultInterval := 1
	defaultOffset := 0
	for _, f := range frequencies {
		t.Run("should have default interval and offset values", func(t *testing.T) {
			if f.Interval() != defaultInterval {
				t.Errorf("Frequency %v should have default interval %v, but got %v", f, defaultInterval, f.Interval())
			}
			if f.Offset() != defaultOffset {
				t.Errorf("Frequency %v should have default offset %v, but got %v", f, defaultOffset, f.Offset())
			}
		})
	}
}

func TestOffsetInterval(t *testing.T) {
	defaultInterval := 1
	defaultOffset := 0

	type values struct {
		interval int
		offset   int
	}
	type errors struct {
		interval bool
		offset   bool
	}
	tests := []struct {
		name    string
		args    values
		want    values
		wantErr errors
	}{
		{
			name:    "should set interval and offset",
			args:    values{interval: 2, offset: 1},
			want:    values{interval: 2, offset: 1},
			wantErr: errors{interval: false, offset: false},
		},
		{
			name:    "should set interval and offset",
			args:    values{interval: 1000, offset: 1000},
			want:    values{interval: 1000, offset: 1000},
			wantErr: errors{interval: false, offset: false},
		},
		{
			name:    "should return error when setting interval",
			args:    values{interval: 0, offset: 1},
			want:    values{interval: defaultInterval, offset: 1},
			wantErr: errors{interval: true, offset: false},
		},
		{
			name:    "should return error when setting offset",
			args:    values{interval: 2, offset: -1},
			want:    values{interval: 2, offset: defaultOffset},
			wantErr: errors{interval: false, offset: true},
		},
		{
			name:    "should return error when setting both interval and offset",
			args:    values{interval: 0, offset: -1},
			want:    values{interval: defaultInterval, offset: defaultOffset},
			wantErr: errors{interval: true, offset: true},
		},
		{
			name:    "should return error when setting both interval and offset",
			args:    values{interval: -1000, offset: -1000},
			want:    values{interval: defaultInterval, offset: defaultOffset},
			wantErr: errors{interval: true, offset: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frequencies := emptyFrequencies(t)

			for _, f := range frequencies {
				err := f.SetInterval(tt.args.interval)
				if (err != nil) != tt.wantErr.interval {
					t.Errorf("SetInterval() error = %v, wantErr %v", err, tt.wantErr.interval)
				}
				if f.Interval() != tt.want.interval {
					t.Errorf("Interval() = %v, want %v", f.Interval(), tt.want.interval)
				}
				err = f.SetOffset(tt.args.offset)
				if (err != nil) != tt.wantErr.offset {
					t.Errorf("SetOffset() error = %v, wantErr %v", err, tt.wantErr.offset)
				}
				if f.Offset() != tt.want.offset {
					t.Errorf("Offset() = %v, want %v", f.Interval(), tt.want.offset)
				}
			}
		})
	}
}

func TestNewHourFrequency(t *testing.T) {
	type args struct {
		atMinutes []int
	}
	tests := []struct {
		name    string
		args    args
		want    Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{atMinutes: []int{0, 1, 59}},
			want:    Frequency{interval: 1, timePeriod: TimePeriodHour, atMinutes: []int{0, 1, 59}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{atMinutes: []int{0, 1, 60}},
			want:    Frequency{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHourFrequency(tt.args.atMinutes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHourFrequency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHourFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDayFrequency(t *testing.T) {
	type args struct {
		atMinutes []int
		atHours   []int
	}
	tests := []struct {
		name    string
		args    args
		want    Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}},
			want:    Frequency{interval: 1, timePeriod: TimePeriodDay, atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{atMinutes: []int{0, 1, 60}, atHours: []int{0, 1, 10, 22, 23}},
			want:    Frequency{},
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 23, 24}},
			want:    Frequency{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDayFrequency(tt.args.atMinutes, tt.args.atHours)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDayFrequency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDayFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWeekFrequency(t *testing.T) {
	type args struct {
		atMinutes []int
		atHours   []int
		onDays    []time.Weekday
	}
	tests := []struct {
		name    string
		args    args
		want    Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}, onDays: []time.Weekday{time.Sunday, time.Monday, time.Saturday}},
			want:    Frequency{interval: 1, timePeriod: TimePeriodWeek, atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}, onDaysOfWeek: []time.Weekday{time.Sunday, time.Monday, time.Saturday}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{atMinutes: []int{0, 1, 60}, atHours: []int{0, 1, 10, 22, 23}, onDays: []time.Weekday{time.Sunday, time.Monday, time.Saturday}},
			want:    Frequency{},
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 23, 24}, onDays: []time.Weekday{time.Sunday, time.Monday, time.Saturday}},
			want:    Frequency{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWeekFrequency(tt.args.atMinutes, tt.args.atHours, tt.args.onDays)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWeekFrequency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWeekFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMonthFrequency(t *testing.T) {
	type args struct {
		atMinutes []int
		atHours   []int
		onDays    []int
	}
	tests := []struct {
		name    string
		args    args
		want    Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}, onDays: []int{1, 2, 30, 31}},
			want:    Frequency{interval: 1, timePeriod: TimePeriodMonth, atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 22, 23}, onDaysOfMonth: []int{1, 2, 30, 31}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{atMinutes: []int{0, 1, 60}, atHours: []int{0, 1, 10, 22, 23}, onDays: []int{1, 2, 30, 31}},
			want:    Frequency{},
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 23, 24}, onDays: []int{1, 2, 30, 31}},
			want:    Frequency{},
			wantErr: true,
		},
		{
			name:    "should return error with days < 1",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 23, 24}, onDays: []int{0, 1, 2, 30, 31}},
			want:    Frequency{},
			wantErr: true,
		},
		{
			name:    "should return error with days > 31",
			args:    args{atMinutes: []int{0, 1, 59}, atHours: []int{0, 1, 10, 23, 24}, onDays: []int{1, 2, 30, 31, 32}},
			want:    Frequency{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMonthFrequency(tt.args.atMinutes, tt.args.atHours, tt.args.onDays)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMonthFrequency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMonthFrequency() = %v, want %v", got, tt.want)
			}
		})
	}
}
