package schedule

import (
	"fmt"
)

// Frequency defines how often an event occurs
type Frequency struct {
	interval      uint8
	timePeriod    TimePeriod
	atMinutes     []uint8
	atHours       []uint8
	onDaysOfWeek  []Day
	onDaysOfMonth []uint8
}

// NewHourFrequency creates a new struct that represents an hour frequency
func NewHourFrequency(interval uint8, atMinutes []uint8) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:   interval,
		timePeriod: TimePeriodHour,
		atMinutes:  atMinutes,
	}, nil
}

// NewDayFrequency creates a new struct that represents a day frequency
func NewDayFrequency(interval uint8, atMinutes []uint8, atHours []uint8) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}
	if err := validateHours(atHours); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:   interval,
		timePeriod: TimePeriodDay,
		atMinutes:  atMinutes,
		atHours:    atHours,
	}, nil
}

// NewWeekFrequency creates a new struct that represents a week frequency
func NewWeekFrequency(interval uint8, atMinutes []uint8, atHours []uint8, onDays []Day) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}
	if err := validateHours(atHours); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:     interval,
		timePeriod:   TimePeriodWeek,
		atMinutes:    atMinutes,
		atHours:      atHours,
		onDaysOfWeek: onDays,
	}, nil
}

// NewMonthFrequency creates a new struct that represents a month frequency
func NewMonthFrequency(interval uint8, atMinutes []uint8, atHours []uint8, onDays []uint8) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}
	if err := validateHours(atHours); err != nil {
		return nil, err
	}
	if err := validateDaysOfMonth(onDays); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:      interval,
		timePeriod:    TimePeriodMonth,
		atMinutes:     atMinutes,
		atHours:       atHours,
		onDaysOfMonth: onDays,
	}, nil
}

func validateMinutes(mins []uint8) error {
	for _, min := range mins {
		if min >= 60 {
			return fmt.Errorf("Minutes must be between 0 and 59, inclusive")
		}
	}
	return nil
}

func validateHours(hrs []uint8) error {
	for _, hr := range hrs {
		if hr >= 24 {
			return fmt.Errorf("Hours must be between 0 and 23, inclusive")
		}
	}
	return nil
}

func validateDaysOfMonth(days []uint8) error {
	for _, day := range days {
		if day < 1 || day > 31 {
			return fmt.Errorf("Days of month must be between 1 and 31, inclusive")
		}
	}
	return nil
}
