package json

import (
	"encoding/json"
	"io"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
)

// Parser handles JSON parsing
type Parser struct {
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
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
	return task.New(at.Name, at.Description, user.ID{})
}
