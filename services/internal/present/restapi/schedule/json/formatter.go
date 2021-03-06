package json

import (
	"encoding/json"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
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

type outSchedule struct {
	ID            usecase.ScheduleID `json:"id"`
	Frequency     string             `json:"frequency"`
	Interval      int                `json:"interval"`
	Offset        int                `json:"offset"`
	AtMinutes     []int              `json:"atMinutes,omitempty"`
	AtHours       []int              `json:"atHours,omitempty"`
	OnDaysOfWeek  []format.Weekday   `json:"onDaysOfWeek,omitempty"`
	OnDaysOfMonth []int              `json:"onDaysOfMonth,omitempty"`
	Paused        bool               `json:"paused"`
	Tasks         []outRecurringTask `json:"tasks"`
}

type outRecurringTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type outTaskID struct {
	ID usecase.TaskID `json:"id"`
}

type outScheduleID struct {
	ID usecase.ScheduleID `json:"id"`
}

type outClearedCompleted struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

// ScheduleID formats a ScheduleID to JSON
func (f *Formatter) ScheduleID(id usecase.ScheduleID) ([]byte, error) {
	o := &outScheduleID{
		ID: id,
	}
	return json.Marshal(o)
}

func scheduleToOut(id usecase.ScheduleID, s *schedule.Schedule) *outSchedule {
	f := s.Frequency()
	outS := outSchedule{
		ID:        id,
		Frequency: f.TimePeriod().String(),
		Interval:  f.Interval(),
		Offset:    f.Offset(),
		Paused:    s.Paused(),
		Tasks:     []outRecurringTask{},
	}
	switch f.TimePeriod() {
	case schedule.TimePeriodHour:
		outS.AtMinutes = f.AtMinutes()
	case schedule.TimePeriodDay:
		outS.AtHours = f.AtHours()
		outS.AtMinutes = f.AtMinutes()
	case schedule.TimePeriodWeek:
		var onDaysOfWeek []format.Weekday
		for _, day := range f.OnDaysOfWeek() {
			onDaysOfWeek = append(onDaysOfWeek, format.Weekday(day))
		}
		outS.OnDaysOfWeek = onDaysOfWeek
		outS.AtHours = f.AtHours()
		outS.AtMinutes = f.AtMinutes()
	case schedule.TimePeriodMonth:
		outS.OnDaysOfMonth = f.OnDaysOfMonth()
		outS.AtHours = f.AtHours()
		outS.AtMinutes = f.AtMinutes()
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
