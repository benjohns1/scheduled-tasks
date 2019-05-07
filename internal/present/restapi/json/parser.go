package json

import (
	"encoding/json"
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
	Paused bool `json:"paused"`
}

func parseAddSchedule(as *addSchedule) (*schedule.Schedule, error) {

	// @TODO: parse actual schedule data
	f, err := schedule.NewHourFrequency([]int{0})
	if err != nil {
		return nil, err
	}
	s := schedule.New(f)
	if as.Paused {
		s.Pause()
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
