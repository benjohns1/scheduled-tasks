package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/core/task"
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

type outSchedule struct {
	ID        usecase.ScheduleID `json:"id"`
	Frequency string             `json:"frequency"`
	AtMinutes []int              `json:"atMinutes"`
	Paused    bool               `json:"paused"`
	Tasks     []outRecurringTask `json:"tasks"`
}

type outRecurringTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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

type outTaskID struct {
	ID usecase.TaskID `json:"id"`
}

type outScheduleID struct {
	ID usecase.ScheduleID `json:"id"`
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
	o := &outTaskID{
		ID: id,
	}
	return json.Marshal(o)
}

// ScheduleID formats a ScheduleID to JSON
func (f *Formatter) ScheduleID(id usecase.ScheduleID) ([]byte, error) {
	o := &outScheduleID{
		ID: id,
	}
	return json.Marshal(o)
}

func taskToOut(id usecase.TaskID, t *task.Task) *outTask {
	return &outTask{
		ID:            id,
		Name:          t.Name(),
		Description:   t.Description(),
		CompletedTime: outTime{t.CompletedTime()},
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
func scheduleToOut(id usecase.ScheduleID, s *schedule.Schedule) *outSchedule {
	outS := outSchedule{
		ID:     id,
		Paused: s.Paused(),
		Tasks:  []outRecurringTask{},
	}
	f := s.Frequency()
	switch f.TimePeriod() {
	case schedule.TimePeriodHour:
		outS.Frequency = "hourly"
		outS.AtMinutes = f.AtMinutes()
		break
	}
	for _, rt := range s.Tasks() {
		oRt := outRecurringTask{Name: rt.Name(), Description: rt.Description()}
		outS.Tasks = append(outS.Tasks, oRt)
	}
	return &outS
}

// Schedule formats a Schedule to JSON
func (f *Formatter) Schedule(sd *usecase.ScheduleData) ([]byte, error) {
	return json.Marshal(scheduleToOut(sd.ScheduleID, sd.Schedule))
}

// ScheduleMap formats a map of Schedules to JSON
func (f *Formatter) ScheduleMap(ss map[usecase.ScheduleID]*schedule.Schedule) ([]byte, error) {
	o := make(map[usecase.ScheduleID]*outSchedule)
	for id, s := range ss {
		o[id] = scheduleToOut(id, s)
	}

	return json.Marshal(o)
}

// Errorf formats a format string and args to JSON
func (f *Formatter) Errorf(format string, a ...interface{}) []byte {
	return f.Error(fmt.Sprintf(format, a...))
}

// Error formats an error message to JSON
func (f *Formatter) Error(a interface{}) []byte {
	outError := &outError{
		Error: fmt.Sprint(a),
	}

	o, mErr := json.Marshal(outError)
	if mErr != nil {
		f.l.Printf("problem marshalling JSON error response: %v (error struct: %v)", mErr, outError)
	}
	return o
}
