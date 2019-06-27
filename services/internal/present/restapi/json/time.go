package json

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Time wraps time.Time for formatting
type Time time.Time

// MarshalJSON formats a time field
func (ft *Time) MarshalJSON() ([]byte, error) {
	var timeStr string
	t := time.Time(*ft)
	if t.IsZero() {
		return []byte("null"), nil
	}
	timeStr = t.Format(OutTimeFormat)
	return []byte(fmt.Sprintf("\"%s\"", timeStr)), nil
}

// Weekday wraps time.Weekday for formatting
type Weekday time.Weekday

// UnmarshalJSON parses a weekday field
func (w *Weekday) UnmarshalJSON(b []byte) error {
	var dayStr string
	if err := json.Unmarshal(b, &dayStr); err != nil {
		return err
	}
	switch strings.ToLower(dayStr) {
	default:
		return fmt.Errorf("unknown day of the week '%v'", dayStr)
	case "sunday":
		*w = Weekday(time.Sunday)
	case "monday":
		*w = Weekday(time.Monday)
	case "tuesday":
		*w = Weekday(time.Tuesday)
	case "wednesday":
		*w = Weekday(time.Wednesday)
	case "thursday":
		*w = Weekday(time.Thursday)
	case "friday":
		*w = Weekday(time.Friday)
	case "saturday":
		*w = Weekday(time.Saturday)
	}
	return nil
}

// MarshalJSON formats a weekday field
func (w *Weekday) MarshalJSON() ([]byte, error) {
	dayStr := time.Weekday(*w).String()
	return []byte(fmt.Sprintf("\"%s\"", dayStr)), nil
}
