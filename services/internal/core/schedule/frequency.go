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

// Offset returns the frequency's offset value
func (f *Frequency) Offset() int {
	return f.offset
}

// Interval returns the frequency's interval value
func (f *Frequency) Interval() int {
	return f.interval
}

// TimePeriod returns the frequency's timePeriod value
func (f *Frequency) TimePeriod() TimePeriod {
	return f.timePeriod
}

// AtMinutes returns the frequency's atMinutes value
func (f *Frequency) AtMinutes() []int {
	return f.atMinutes
}

// AtHours returns the frequency's atHours value
func (f *Frequency) AtHours() []int {
	return f.atHours
}

// OnDaysOfWeek returns the frequency's onDaysOfWeek value
func (f *Frequency) OnDaysOfWeek() []time.Weekday {
	return f.onDaysOfWeek
}

// OnDaysOfMonth returns the frequency's onDaysOfMonth value
func (f *Frequency) OnDaysOfMonth() []int {
	return f.onDaysOfMonth
}

// NewRawFrequency creates a new frequency struct from raw data
func NewRawFrequency(offset int, interval int, timePeriod TimePeriod, atMinutes []int, atHours []int, onDaysOfWeek []time.Weekday, onDaysOfMonth []int) (Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return Frequency{}, err
	}
	if err := validateHours(atHours); err != nil {
		return Frequency{}, err
	}
	if err := validateDaysOfMonth(onDaysOfMonth); err != nil {
		return Frequency{}, err
	}

	return Frequency{offset, interval, timePeriod, atMinutes, atHours, onDaysOfWeek, onDaysOfMonth}, nil
}

// NewHourFrequency creates a new struct that represents an hour frequency
func NewHourFrequency(atMinutes []int) (Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return Frequency{}, err
	}

	return Frequency{
		interval:   1,
		timePeriod: TimePeriodHour,
		atMinutes:  atMinutes,
	}, nil
}

// NewDayFrequency creates a new struct that represents a day frequency
func NewDayFrequency(atMinutes []int, atHours []int) (Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return Frequency{}, err
	}
	if err := validateHours(atHours); err != nil {
		return Frequency{}, err
	}

	return Frequency{
		interval:   1,
		timePeriod: TimePeriodDay,
		atMinutes:  atMinutes,
		atHours:    atHours,
	}, nil
}

// NewWeekFrequency creates a new struct that represents a week frequency
func NewWeekFrequency(atMinutes []int, atHours []int, onDays []time.Weekday) (Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return Frequency{}, err
	}
	if err := validateHours(atHours); err != nil {
		return Frequency{}, err
	}

	return Frequency{
		interval:     1,
		timePeriod:   TimePeriodWeek,
		atMinutes:    atMinutes,
		atHours:      atHours,
		onDaysOfWeek: onDays,
	}, nil
}

// NewMonthFrequency creates a new struct that represents a month frequency
func NewMonthFrequency(atMinutes []int, atHours []int, onDays []int) (Frequency, error) {
	if err := validateMinutes(atMinutes); err != nil {
		return Frequency{}, err
	}
	if err := validateHours(atHours); err != nil {
		return Frequency{}, err
	}
	if err := validateDaysOfMonth(onDays); err != nil {
		return Frequency{}, err
	}

	return Frequency{
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
		return f.calcHourTimes(start, &end)
	case TimePeriodDay:
		return f.calcDayTimes(start, &end)
	case TimePeriodWeek:
		return f.calcWeekTimes(start, &end)
	}
	return nil, fmt.Errorf("timePeriod %v not implemented yet", f.timePeriod)
}

func (f *Frequency) next(after time.Time) (time.Time, error) {
	switch f.timePeriod {
	case TimePeriodNone:
		return time.Time{}, nil
	case TimePeriodHour:
		return getNextTime(f.calcHourTimes(after, nil))
	case TimePeriodDay:
		return getNextTime(f.calcDayTimes(after, nil))
	case TimePeriodWeek:
		return getNextTime(f.calcWeekTimes(after, nil))
	}
	return time.Time{}, fmt.Errorf("timePeriod %v not implemented yet", f.timePeriod)
}

func getNextTime(times []time.Time, err error) (time.Time, error) {
	if err != nil {
		return time.Time{}, err
	}
	if len(times) > 0 {
		return times[0], nil
	}
	return time.Time{}, nil
}

func (f *Frequency) calcWeekTimes(start time.Time, end *time.Time) ([]time.Time, error) {
	/*
		var max int
		if end == nil {
			max = 52
		} else {
			max = (int(end.Sub(start).Hours()/(24*7)) / f.interval) + 1
		}
		times := []time.Time{}
		if len(f.onDaysOfWeek) == 0 {
			return times, nil
		}
		dayOffset = f.offset * 7
	*/

	// Calculate first week

	// @TODO: finish week times calculation
	return []time.Time{}, fmt.Errorf("NOT IMPLEMENTED")
}

func (f *Frequency) calcDayTimes(start time.Time, end *time.Time) ([]time.Time, error) {

	var max int
	if end == nil {
		max = 365
	} else {
		max = (int(end.Sub(start).Hours()/24) / f.interval) + 1
	}
	times := []time.Time{}

	// Calculate first day
	day := start.YearDay() + ((start.YearDay() - 1) % f.interval) + f.offset

	// Add times to the array
	for d := 0; d <= max; d++ {
		for _, hour := range f.atHours {
			for _, min := range f.atMinutes {
				t := time.Date(start.Year(), 1, day, hour, min, 0, 0, start.Location())
				if t.Before(start) {
					continue
				}
				if end == nil {
					return append(times, t), nil
				}
				if t.After(*end) {
					return times, nil
				}
				times = append(times, t)
			}
		}
		day += f.interval
	}
	return times, nil
}

func (f *Frequency) calcHourTimes(start time.Time, end *time.Time) ([]time.Time, error) {

	var max int
	if end == nil {
		max = 365 * 24 // search up to 1 year
	} else {
		max = (int(end.Sub(start).Hours()) / f.interval) + 1
	}
	times := []time.Time{}

	// Calculate first hour
	hour := start.Hour() + (start.Hour() % f.interval) + f.offset

	// Add times to the array
	for h := 0; h <= max; h++ {
		for _, min := range f.atMinutes {
			t := time.Date(start.Year(), start.Month(), start.Day(), hour, min, 0, 0, start.Location())
			if t.Before(start) {
				continue
			}
			if end == nil {
				return append(times, t), nil
			}
			if t.After(*end) {
				return times, nil
			}
			times = append(times, t)
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
