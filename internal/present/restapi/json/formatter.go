package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter formats application data into JSON for output
type Formatter struct {
	l Logger
}

// NewFormatter creates a new Formatter instance
func NewFormatter(l Logger) *Formatter {
	return &Formatter{l: l}
}

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
		return []byte("null"), nil
	}
	timeStr = t.t.Format(outTimeFormat)
	return []byte(fmt.Sprintf("\"%s\"", timeStr)), nil
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

// WriteResponse writes the complete output response
func (f *Formatter) WriteResponse(w http.ResponseWriter, res []byte, statusCode int) {
	f.WriteEmpty(w, statusCode)
	w.Write(res)
}

// WriteEmpty writes a complete empty output response
func (f *Formatter) WriteEmpty(w http.ResponseWriter, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
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

// TaskMap formats a map of Tasks to JSON
func (f *Formatter) TaskMap(ts map[usecase.TaskID]*core.Task) ([]byte, error) {
	o := make(map[usecase.TaskID]*outTask)
	for id, t := range ts {
		o[id] = taskToOut(id, t)
	}

	return json.Marshal(o)
}

// Error formats an error message to JSON
func (f *Formatter) Error(err error) []byte {
	outError := &outError{
		Error: err.Error(),
	}

	o, mErr := json.Marshal(outError)
	if mErr != nil {
		f.l.Printf("problem marshalling JSON error response: %v (error struct: %v)", mErr, outError)
	}
	return o
}
