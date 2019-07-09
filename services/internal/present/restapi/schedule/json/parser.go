package json

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	parse "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
)

// Parser handles JSON parsing
type Parser struct {
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// AddSchedule parses addSchedule request JSON data into a core Schedule struct
func (p *Parser) AddSchedule(b io.Reader) (*schedule.Schedule, error) {
	var addSchedule addSchedule
	err := json.NewDecoder(b).Decode(&addSchedule)
	if err != nil {
		return nil, err
	}
	return parseAddSchedule(&addSchedule)
}

type addSchedule struct {
	Frequency     string             `json:"frequency"`
	Interval      *int               `json:"interval"`
	Offset        *int               `json:"offset"`
	AtMinutes     []int              `json:"atMinutes"`
	AtHours       []int              `json:"atHours"`
	OnDaysOfWeek  []parse.Weekday    `json:"onDaysOfWeek"`
	OnDaysOfMonth []int              `json:"onDaysOfMonth"`
	Paused        bool               `json:"paused"`
	Tasks         []addRecurringTask `json:"tasks"`
}

func parseAddSchedule(as *addSchedule) (*schedule.Schedule, error) {

	var f schedule.Frequency
	var err error
	switch as.Frequency {
	case "Hour":
		f, err = schedule.NewHourFrequency(as.AtMinutes)
	case "Day":
		f, err = schedule.NewDayFrequency(as.AtMinutes, as.AtHours)
	case "Week":
		var onDaysOfWeek []time.Weekday
		for _, d := range as.OnDaysOfWeek {
			onDaysOfWeek = append(onDaysOfWeek, time.Weekday(d))
		}
		f, err = schedule.NewWeekFrequency(as.AtMinutes, as.AtHours, onDaysOfWeek)
	case "Month":
		f, err = schedule.NewMonthFrequency(as.AtMinutes, as.AtHours, as.OnDaysOfMonth)
	default:
		return nil, fmt.Errorf("invalid frequency '%v', should be 'Hour', 'Day', 'Week', or 'Month'", as.Frequency)
	}
	if err != nil {
		return nil, err
	}
	if as.Interval != nil {
		err = f.SetInterval(*as.Interval)
		if err != nil {
			return nil, err
		}
	}
	if as.Offset != nil {
		err = f.SetOffset(*as.Offset)
		if err != nil {
			return nil, err
		}
	}

	s := schedule.New(f, user.ID{})
	if as.Paused {
		s.Pause()
	}
	for _, rt := range as.Tasks {
		s.AddTask(schedule.NewRecurringTask(rt.Name, rt.Description))
	}
	return s, nil
}

// AddRecurringTask parses addRecurringTask request JSON into a core RecurringTask struct
func (p *Parser) AddRecurringTask(b io.Reader) (schedule.RecurringTask, error) {
	var addRecurringTask addRecurringTask
	if err := json.NewDecoder(b).Decode(&addRecurringTask); err != nil {
		return schedule.RecurringTask{}, err
	}
	return parseAddRecurringTask(&addRecurringTask), nil
}

type addRecurringTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func parseAddRecurringTask(rt *addRecurringTask) schedule.RecurringTask {
	return schedule.NewRecurringTask(rt.Name, rt.Description)
}
