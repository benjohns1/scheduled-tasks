package schedule

import (
	"reflect"
	"testing"
	"time"
)

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
