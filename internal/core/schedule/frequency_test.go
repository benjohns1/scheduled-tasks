package schedule

import (
	"reflect"
	"testing"
)

func TestNewHourFrequency(t *testing.T) {
	type args struct {
		interval  uint8
		atMinutes []uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}},
			want:    &Frequency{interval: 0, timePeriod: TimePeriodHour, atMinutes: []uint8{0, 1, 59}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 60}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHourFrequency(tt.args.interval, tt.args.atMinutes)
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
		interval  uint8
		atMinutes []uint8
		atHours   []uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}},
			want:    &Frequency{interval: 0, timePeriod: TimePeriodDay, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 60}, atHours: []uint8{0, 1, 10, 22, 23}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 23, 24}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDayFrequency(tt.args.interval, tt.args.atMinutes, tt.args.atHours)
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
		interval  uint8
		atMinutes []uint8
		atHours   []uint8
		onDays    []Day
	}
	tests := []struct {
		name    string
		args    args
		want    *Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}, onDays: []Day{Sunday, Monday, Saturday}},
			want:    &Frequency{interval: 0, timePeriod: TimePeriodWeek, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}, onDaysOfWeek: []Day{Sunday, Monday, Saturday}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 60}, atHours: []uint8{0, 1, 10, 22, 23}, onDays: []Day{Sunday, Monday, Saturday}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 23, 24}, onDays: []Day{Sunday, Monday, Saturday}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWeekFrequency(tt.args.interval, tt.args.atMinutes, tt.args.atHours, tt.args.onDays)
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
		interval  uint8
		atMinutes []uint8
		atHours   []uint8
		onDays    []uint8
	}
	tests := []struct {
		name    string
		args    args
		want    *Frequency
		wantErr bool
	}{
		{
			name:    "should return a valid struct",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}, onDays: []uint8{1, 2, 30, 31}},
			want:    &Frequency{interval: 0, timePeriod: TimePeriodMonth, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 22, 23}, onDaysOfMonth: []uint8{1, 2, 30, 31}},
			wantErr: false,
		},
		{
			name:    "should return error with minutes >= 60",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 60}, atHours: []uint8{0, 1, 10, 22, 23}, onDays: []uint8{1, 2, 30, 31}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with hours >= 24",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 23, 24}, onDays: []uint8{1, 2, 30, 31}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with days < 1",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 23, 24}, onDays: []uint8{0, 1, 2, 30, 31}},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "should return error with days > 31",
			args:    args{interval: 0, atMinutes: []uint8{0, 1, 59}, atHours: []uint8{0, 1, 10, 23, 24}, onDays: []uint8{1, 2, 30, 31, 32}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMonthFrequency(tt.args.interval, tt.args.atMinutes, tt.args.atHours, tt.args.onDays)
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

func Test_validateMinutes(t *testing.T) {
	type args struct {
		mins []uint8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMinutes(tt.args.mins); (err != nil) != tt.wantErr {
				t.Errorf("validateMinutes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateHours(t *testing.T) {
	type args struct {
		hrs []uint8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateHours(tt.args.hrs); (err != nil) != tt.wantErr {
				t.Errorf("validateHours() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateDaysOfMonth(t *testing.T) {
	type args struct {
		days []uint8
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDaysOfMonth(tt.args.days); (err != nil) != tt.wantErr {
				t.Errorf("validateDaysOfMonth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
