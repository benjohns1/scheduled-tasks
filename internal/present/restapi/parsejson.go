package restapi

import (
	"encoding/json"
	"github.com/benjohns1/scheduled-tasks/internal/core"
	"io"
)

type addTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func addTaskToTask(at *addTask) *core.Task {
	return core.NewTask(at.Name, at.Description)
}

func addTaskFromJSON(b io.ReadCloser) (*core.Task, error) {
	var addTask addTask
	err := json.NewDecoder(b).Decode(&addTask)
	if err != nil {
		return nil, err
	}
	return addTaskToTask(&addTask), nil
}
