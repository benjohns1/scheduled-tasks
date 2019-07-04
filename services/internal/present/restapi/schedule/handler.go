package schedule

import (
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
	responseMapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	mapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/schedule/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter defines the formatter interface for output responses
type Formatter interface {
	WriteResponse(w http.ResponseWriter, res []byte, statusCode int)
	WriteEmpty(w http.ResponseWriter, statusCode int)
	Errorf(format string, a ...interface{}) []byte
	Error(a interface{}) []byte
	Schedule(sd *usecase.ScheduleData) ([]byte, error)
	ScheduleID(id usecase.ScheduleID) ([]byte, error)
	ScheduleMap(ss map[usecase.ScheduleID]*schedule.Schedule) ([]byte, error)
}

// Parser defines the parser interface for parsing input requests
type Parser interface {
	AddSchedule(b io.Reader) (*schedule.Schedule, error)
	AddRecurringTask(b io.Reader) (schedule.RecurringTask, error)
}

// Handle adds schedule handling endpoints
func Handle(r *httprouter.Router, prefix string, l Logger, rf responseMapper.ResponseFormatter, checkSchedule chan<- bool, scheduleRepo usecase.ScheduleRepo) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(rf)

	sPre := prefix + "/schedule"
	r.GET(sPre+"/", auth.Handler(listSchedules(l, f, scheduleRepo)))
	r.GET(sPre+"/:scheduleID", auth.Handler(getSchedule(l, f, scheduleRepo)))
	r.DELETE(sPre+"/:scheduleID", auth.Handler(removeSchedule(l, f, checkSchedule, scheduleRepo)))
	r.POST(sPre+"/", auth.Handler(addSchedule(l, f, p, checkSchedule, scheduleRepo)))
	r.PUT(sPre+"/:scheduleID/pause", auth.Handler(pauseSchedule(l, f, checkSchedule, scheduleRepo)))
	r.PUT(sPre+"/:scheduleID/unpause", auth.Handler(unpauseSchedule(l, f, checkSchedule, scheduleRepo)))

	rtPre := sPre + "/:scheduleID/task"
	r.POST(rtPre+"/", auth.Handler(addRecurringTask(l, f, p, scheduleRepo)))
}

func listSchedules(l Logger, f Formatter, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ss, err := usecase.ListSchedules(scheduleRepo)
		if err != nil {
			l.Printf("error retrieving schedule list: %v", err)
			f.WriteResponse(w, f.Error("Error: couldn't retrieve schedules"), 500)
			return
		}

		o, e := f.ScheduleMap(ss)

		if e != nil {
			l.Printf("error encoding schedule map: %v", e)
			f.WriteResponse(w, f.Error("Error encoding schedule data"), 500)
		}
		f.WriteResponse(w, o, 200)
	}
}

func addSchedule(l Logger, f Formatter, p Parser, checkSchedule chan<- bool, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s, err := p.AddSchedule(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Printf("error parsing addSchedule data: %v", err)
			f.WriteResponse(w, f.Errorf("Error: could not parse schedule data: %v", err), 400)
			return
		}
		sID, ucerr := usecase.AddSchedule(scheduleRepo, s, checkSchedule)
		if ucerr != nil {
			l.Printf("error adding schedule: %v", ucerr)
			f.WriteResponse(w, f.Error("Error: could not add schedule data"), 500)
			return
		}
		o, err := f.ScheduleID(sID)
		if err != nil {
			f.WriteResponse(w, f.Error("Schedule created, but there was an error formatting the response Schedule ID"), 201)
			return
		}
		f.WriteResponse(w, o, 201)
	}
}

func getSchedule(l Logger, f Formatter, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(params.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		sd, ucerr := usecase.GetSchedule(scheduleRepo, id)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Schedule ID %d not found", id), 404)
				return
			}
			l.Printf("error retrieving schedule ID %d: %v", id, ucerr)
			f.WriteResponse(w, f.Errorf("Error: couldn't retrieve schedule ID %d", id), 500)
			return
		}

		o, err := f.Schedule(sd)
		if err != nil {
			l.Printf("error encoding schedule map: %v", err)
			f.WriteResponse(w, f.Error("Error encoding schedule data"), 500)
		}
		f.WriteResponse(w, o, 200)
	}
}

func removeSchedule(l Logger, f Formatter, checkSchedule chan<- bool, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		ucerr := usecase.RemoveSchedule(scheduleRepo, id, checkSchedule)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Schedule ID %d not found", id), 404)
				return
			}
			l.Printf("error removing schedule: %v", ucerr)
			f.WriteResponse(w, f.Error("Error removing schedule"), 500)
			return
		}
		f.WriteEmpty(w, 204)
	}
}

func pauseSchedule(l Logger, f Formatter, checkSchedule chan<- bool, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		ucerr := usecase.PauseSchedule(scheduleRepo, id, checkSchedule)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Schedule ID %d not found", id), 404)
				return
			}
			l.Printf("error pausing schedule: %v", ucerr)
			f.WriteResponse(w, f.Error("Error pausing schedule"), 500)
			return
		}
		f.WriteEmpty(w, 204)
	}
}

func unpauseSchedule(l Logger, f Formatter, checkSchedule chan<- bool, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		ucerr := usecase.UnpauseSchedule(scheduleRepo, id, checkSchedule)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Schedule ID %d not found", id), 404)
				return
			}
			l.Printf("error unpausing schedule: %v", ucerr)
			f.WriteResponse(w, f.Error("Error unpausing schedule"), 500)
			return
		}
		f.WriteEmpty(w, 204)
	}
}

func addRecurringTask(l Logger, f Formatter, p Parser, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// Get schedule ID
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)

		// Parse recurring task data
		rt, err := p.AddRecurringTask(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Printf("error parsing addRecurringTask data: %v", err)
			f.WriteResponse(w, f.Errorf("Error: could not parse recurring task data: %v", err), 400)
			return
		}

		// Add recurring task
		ucerr := usecase.AddRecurringTask(scheduleRepo, id, rt)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Schedule ID %d not found", id), 404)
				return
			}
			if ucerr.Code() == usecase.ErrDuplicateRecord {
				f.WriteResponse(w, f.Errorf("Recurring task already exists for this schedule, can't add duplicate tasks with the same data"), 400)
				return
			}
			l.Printf("error adding task to schedule: %v", ucerr)
			f.WriteResponse(w, f.Error("Error adding task to schedule"), 500)
			return
		}
		f.WriteEmpty(w, 201)
	}
}
