package json

import (
	"encoding/json"
	"github.com/benjohns1/scheduled-tasks/internal/core"
)

// Parser handles JSON parsing
type Parser struct {
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// AddTask parses addTask request JSON data into a core Task struct
func (p *Parser) AddTask(b Reader) (*core.Task, error) {
	var addTask addTask
	err := json.NewDecoder(b).Decode(&addTask)
	if err != nil {
		return nil, err
	}
	return addTaskToTask(&addTask), nil
}

type addTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func addTaskToTask(at *addTask) *core.Task {
	return core.NewTask(at.Name, at.Description)
}

// Reader defines io.Reader interface needed to parse from JSON
type Reader interface {
	Read(p []byte) (n int, err error)
}
