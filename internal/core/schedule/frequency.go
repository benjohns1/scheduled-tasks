package schedule

import (
	"fmt"
	"time"
)

// Frequency defines how often an event occurs
type Frequency struct {
	offset        int
	interval      int
	timePeriod    TimePeriod
	atMinutes     []int
	atHours       []int
	onDaysOfWeek  []time.Weekday
	onDaysOfMonth []int
}

// NewHourFrequency creates a new struct that represents an hour frequency
func NewHourFrequency(atMinutes []int) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:   1,
		timePeriod: TimePeriodHour,
		atMinutes:  atMinutes,
	}, nil
}

// NewDayFrequency creates a new struct that represents a day frequency
func NewDayFrequency(atMinutes []int, atHours []int) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}
	if err := validateHours(atHours); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:   1,
		timePeriod: TimePeriodDay,
		atMinutes:  atMinutes,
		atHours:    atHours,
	}, nil
}

// NewWeekFrequency creates a new struct that represents a week frequency
func NewWeekFrequency(atMinutes []int, atHours []int, onDays []time.Weekday) (*Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return nil, err
	}
	if err := validateHours(atHours); err != nil {
		return nil, err
	}

	return &Frequency{
		interval:     1,
		timePeriod:   TimePeriodWeek,
		atMinutes:    atMinutes,
		atHours:      atHours,
		onDaysOfWeek: onDays,
	}, nil
}

// NewMonthFrequency creates a new struct that represents a month frequency
func NewMonthFrequency(atMinutes []int, atHours []int, onDays []int) (*Frequency, error) {
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
		interval:      1,
		timePeriod:    TimePeriodMonth,
		atMinutes:     atMinutes,
		atHours:       atHours,
		onDaysOfMonth: onDays,
	}, nil
}

// SetOffset sets an offset from the base of each starting time
// Base starting times for each time period:
//  - Hourly:  Midnight
//  - Daily:   1st day of the year
//  - Weekly:  1st week of the year
//  - Monthly: January
func (f *Frequency) SetOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset %v must be 0 or greater", offset)
	}
	f.offset = offset
	return nil
}

// SetInterval sets an interval between each time period
func (f *Frequency) SetInterval(interval int) error {
	if interval < 1 {
		return fmt.Errorf("interval %v must be greater than 0", interval)
	}
	f.interval = interval
	return nil
}

// times returns all times between the specified start and end time (inclusive) that occur for this frequency
func (f *Frequency) times(start time.Time, end time.Time) ([]time.Time, error) {
	if end.Before(start) {
		return nil, fmt.Errorf("end time %v is before start time %v", end, start)
	}

	switch f.timePeriod {
	case TimePeriodNone:
		return []time.Time{}, nil
	case TimePeriodHour:
		return f.calcHourTimes(start, end)
	}
	return nil, fmt.Errorf("not implemented, yet")
}

func (f *Frequency) calcHourTimes(start time.Time, end time.Time) ([]time.Time, error) {
	maxHour := (int(end.Sub(start).Hours()) / f.interval) + 1
	times := []time.Time{}

	// Calculate first hour
	hour := start.Hour() + (start.Hour() % f.interval)

	// Add times to the array
	for hri := 0; hri <= maxHour; hri++ {
		for _, min := range f.atMinutes {
			time := time.Date(start.Year(), start.Month(), start.Day(), hour, min, 0, 0, start.Location())
			if time.Before(start) {
				continue
			}
			if time.After(end) {
				return times, nil
			}
			times = append(times, time)
		}
		hour += f.interval
	}
	return times, nil
}

func validateMinutes(mins []int) error {
	for _, min := range mins {
		if min < 0 || min > 59 {
			return fmt.Errorf("Minutes must be between 0 and 59, inclusive")
		}
	}
	return nil
}

func validateHours(hrs []int) error {
	for _, hr := range hrs {
		if hr < 0 || hr > 23 {
			return fmt.Errorf("Hours must be between 0 and 23, inclusive")
		}
	}
	return nil
}

func validateDaysOfMonth(days []int) error {
	for _, day := range days {
		if day < 1 || day > 31 {
			return fmt.Errorf("Days of month must be between 1 and 31, inclusive")
		}
	}
	return nil
}
