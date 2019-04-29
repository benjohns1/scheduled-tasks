package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
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

type outID struct {
	ID usecase.TaskID `json:"id"`
}

type outError struct {
	Error string `json:"error"`
}

type outClearedCompleted struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

func (t *outTime) MarshalJSON() ([]byte, error) {
	var timeStr string
	if t.t.IsZero() {
		return []byte("null"), nil
	}
	timeStr = t.t.Format(outTimeFormat)
	return []byte(fmt.Sprintf("\"%s\"", timeStr)), nil
}

func writeResponse(w http.ResponseWriter, res []byte, statusCode int) {
	writeEmpty(w, statusCode)
	w.Write(res)
}

func writeEmpty(w http.ResponseWriter, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

func clearedCompletedToJSON(count int) ([]byte, error) {
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

func idToJSON(id usecase.TaskID) ([]byte, error) {

	o := &outID{
		ID: id,
	}
	return json.Marshal(o)
}

func taskToOut(id usecase.TaskID, t *core.Task) *outTask {
	return &outTask{
		ID:            id,
		Name:          t.Name(),
		Description:   t.Description(),
		CompletedTime: outTime{t.CompletedTime()},
	}
}

func taskMapToJSON(ts map[usecase.TaskID]*core.Task) ([]byte, error) {
	o := make(map[usecase.TaskID]*outTask)
	for id, t := range ts {
		o[id] = taskToOut(id, t)
	}

	return json.Marshal(o)
}

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

func errorToJSON(l Logger, err error) []byte {
	outError := &outError{
		Error: err.Error(),
	}

	o, mErr := json.Marshal(outError)
	if mErr != nil {
		l.Printf("problem marshalling JSON error response: %v (error struct: %v)", mErr, outError)
	}
	return o
}
