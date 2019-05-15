package json

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/core/task"
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
	Frequency string             `json:"frequency"`
	AtMinutes []int              `json:"atMinutes"`
	Paused    bool               `json:"paused"`
	Tasks     []addRecurringTask `json:"tasks"`
}

func parseAddSchedule(as *addSchedule) (*schedule.Schedule, error) {

	// @TODO: parse actual schedule data
	var f schedule.Frequency
	var err error
	switch as.Frequency {
	case "hourly":
		f, err = schedule.NewHourFrequency(as.AtMinutes)
	default:
		return nil, fmt.Errorf("invalid frequency '%v', should be 'hourly', 'daily', 'weekly', or 'monthly'", as.Frequency)
	}
	if err != nil {
		return nil, err
	}
	s := schedule.New(f)
	if as.Paused {
		s.Pause()
	}
	for _, rt := range as.Tasks {
		s.AddTask(schedule.NewRecurringTask(rt.Name, rt.Description))
	}
	return s, nil
}

// AddTask parses addTask request JSON data into a core Task struct
func (p *Parser) AddTask(b io.Reader) (*task.Task, error) {
	var addTask addTask
	err := json.NewDecoder(b).Decode(&addTask)
	if err != nil {
		return nil, err
	}
	return parseAddTask(&addTask), nil
}

type addTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func parseAddTask(at *addTask) *task.Task {
	return task.New(at.Name, at.Description)
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
