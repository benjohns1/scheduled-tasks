package restapi

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

type outTask struct {
	ID            usecase.TaskID `json:"id"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	CompletedTime outTime        `json:"completedTime"`
}

const outTimeFormat = time.RFC3339Nano

type outTime struct {
	t time.Time
}

func (t *outTime) MarshalJSON() ([]byte, error) {
	var timeStr string
	if t.t.IsZero() {
		timeStr = ""
	} else {
		timeStr = t.t.Format(outTimeFormat)
	}
	return []byte(fmt.Sprintf("\"%s\"", timeStr)), nil
}

func taskToOut(id usecase.TaskID, t *core.Task) *outTask {
	return &outTask{
		ID:            id,
		Name:          t.Name(),
		Description:   t.Description(),
		CompletedTime: outTime{t.CompletedTime()},
	}
}

func taskMapToJSON(ts map[usecase.TaskID]*core.Task) []byte {
	o := make(map[usecase.TaskID]*outTask)
	for id, t := range ts {
		o[id] = taskToOut(id, t)
	}

	data, err := json.Marshal(o)
	if err != nil {
		log.Printf("error marshalling task map: %v", err)
		return errorToJSON(fmt.Errorf("Error parsing task data"))
	}
	return data
}

func errorToJSON(err error) []byte {
	data, err := json.Marshal(err)
	return data
}
