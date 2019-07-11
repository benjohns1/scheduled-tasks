package json

import (
	"encoding/json"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	format "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Formatter formats application data into JSON for output
type Formatter struct {
	format.ResponseFormatter
}

// NewFormatter creates a new Formatter instance
func NewFormatter(rf format.ResponseFormatter) *Formatter {
	return &Formatter{rf}
}

type outTask struct {
	ID            usecase.TaskID `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CompletedTime format.Time    `json:"completedTime"`
	CreatedTime   format.Time    `json:"createdTime"`
}

type outTaskID struct {
	ID usecase.TaskID `json:"id"`
}

type outClearedCompleted struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

// ClearedCompleted formats a count of cleared tasks to JSON
func (f *Formatter) ClearedCompleted(count int) ([]byte, error) {
	var message string
	if count > 0 {
		message = "Cleared all completed tasks"
	} else {
		message = "No completed tasks to clear"
	}
	o := &outClearedCompleted{
		Count:   count,
		Message: message,
	}
	return json.Marshal(o)
}

// TaskID formats a TaskID to JSON
func (f *Formatter) TaskID(id usecase.TaskID) ([]byte, error) {
	o := &outTaskID{
		ID: id,
	}
	return json.Marshal(o)
}

func taskToOut(id usecase.TaskID, t *task.Task) *outTask {
	return &outTask{
		ID:            id,
		Name:          t.Name(),
		Description:   t.Description(),
		CompletedTime: format.Time(t.CompletedTime()),
		CreatedTime:   format.Time(t.CreatedTime()),
	}
}

// Task formats a Task to JSON
func (f *Formatter) Task(td *usecase.TaskData) ([]byte, error) {
	return json.Marshal(taskToOut(td.TaskID, td.Task))
}

// TaskMap formats a map of Tasks to JSON
func (f *Formatter) TaskMap(ts map[usecase.TaskID]*task.Task) ([]byte, error) {
	o := make(map[usecase.TaskID]*outTask)
	for id, t := range ts {
		o[id] = taskToOut(id, t)
	}

	return json.Marshal(o)
}
