package restapi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/internal/core/task"
	mapper "github.com/benjohns1/scheduled-tasks/internal/present/restapi/json"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter defines the formatter interface for output responses
type Formatter interface {
	WriteResponse(w http.ResponseWriter, res []byte, statusCode int)
	WriteEmpty(w http.ResponseWriter, statusCode int)
	ClearedCompleted(count int) ([]byte, error)
	TaskID(id usecase.TaskID) ([]byte, error)
	Task(td *usecase.TaskData) ([]byte, error)
	TaskMap(ts map[usecase.TaskID]*task.Task) ([]byte, error)
	Schedule(sd *usecase.ScheduleData) ([]byte, error)
	ScheduleID(id usecase.ScheduleID) ([]byte, error)
	ScheduleMap(ss map[usecase.ScheduleID]*schedule.Schedule) ([]byte, error)
	Errorf(format string, a ...interface{}) []byte
	Error(a interface{}) []byte
}

// Parser defines the parser interface for parsing input requests
type Parser interface {
	AddTask(b io.Reader) (*task.Task, error)
	AddSchedule(b io.Reader) (*schedule.Schedule, error)
	AddRecurringTask(b io.Reader) (schedule.RecurringTask, error)
}

// New creates a REST API server
func New(l Logger, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo) (api http.Handler) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(l)
	r := httprouter.New()

	prefix := "/api/v1"

	tPre := prefix + "/task"
	r.GET(tPre+"/", listTasks(l, f, taskRepo))
	r.GET(tPre+"/:taskID", getTask(l, f, taskRepo))
	r.POST(tPre+"/", addTask(l, f, p, taskRepo))
	r.PUT(tPre+"/:taskID/complete", completeTask(l, f, taskRepo))
	r.DELETE(tPre+"/:taskID", clearTask(l, f, taskRepo))
	r.POST(tPre+"/clear", clearCompletedTasks(l, f, taskRepo))

	sPre := prefix + "/schedule"
	r.GET(sPre+"/", listSchedules(l, f, scheduleRepo))
	r.GET(sPre+"/:scheduleID", getSchedule(l, f, scheduleRepo))
	r.POST(sPre+"/", addSchedule(l, f, p, scheduleRepo))
	r.PUT(sPre+"/:scheduleID/pause", pauseSchedule(l, f, scheduleRepo))
	r.PUT(sPre+"/:scheduleID/unpause", unpauseSchedule(l, f, scheduleRepo))

	rtPre := sPre + "/:scheduleID/task"
	r.POST(rtPre+"/", addRecurringTask(l, f, p, scheduleRepo))

	return r
}

// Serve starts an API server
func Serve(l Logger, api http.Handler) {

	// Start API server
	port := 8080
	if val, err := strconv.Atoi(os.Getenv("APPLICATION_PORT")); err == nil {
		port = val
	}

	l.Printf("starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), api)
	l.Printf("server exiting")
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

func addSchedule(l Logger, f Formatter, p Parser, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		s, err := p.AddSchedule(r.Body)
		defer r.Body.Close()
		if err != nil {
			l.Printf("error parsing addSchedule data: %v", err)
			f.WriteResponse(w, f.Errorf("Error: could not parse schedule data: %v", err), 400)
			return
		}
		sID, ucerr := usecase.AddSchedule(scheduleRepo, s)
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
			f.WriteResponse(w, f.Error("Error encoding task data"), 500)
		}
		f.WriteResponse(w, o, 200)
	}
}

func pauseSchedule(l Logger, f Formatter, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		ucerr := usecase.PauseSchedule(scheduleRepo, id)
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

func unpauseSchedule(l Logger, f Formatter, scheduleRepo usecase.ScheduleRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		scheduleIDInt, err := strconv.Atoi(ps.ByName("scheduleID"))
		if err != nil {
			l.Printf("valid schedule ID required")
			f.WriteResponse(w, f.Error("Error: valid schedule ID required"), 404)
			return
		}
		id := usecase.ScheduleID(scheduleIDInt)
		ucerr := usecase.UnpauseSchedule(scheduleRepo, id)
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

func listTasks(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ts, ucerr := usecase.ListTasks(taskRepo)
		if ucerr != nil {
			l.Printf("error retrieving task list: %v", ucerr)
			f.WriteResponse(w, f.Error("Error: couldn't retrieve tasks"), 500)
			return
		}
		o, err := f.TaskMap(ts)
		if err != nil {
			l.Printf("error encoding task map: %v", err)
			f.WriteResponse(w, f.Error("Error encoding task data"), 500)
		}
		f.WriteResponse(w, o, 200)
	}
}

func getTask(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		taskIDInt, err := strconv.Atoi(params.ByName("taskID"))
		if err != nil {
			l.Printf("valid task ID required")
			f.WriteResponse(w, f.Error("Error: valid task ID required"), 404)
			return
		}
		id := usecase.TaskID(taskIDInt)
		td, ucerr := usecase.GetTask(taskRepo, id)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Task ID %d not found", id), 404)
				return
			}
			l.Printf("error retrieving task ID %d: %v", id, ucerr)
			f.WriteResponse(w, f.Errorf("Error: couldn't retrieve task ID %d", id), 500)
			return
		}

		o, err := f.Task(td)
		if err != nil {
			l.Printf("error encoding task map: %v", err)
			f.WriteResponse(w, f.Error("Error encoding task data"), 500)
		}
		f.WriteResponse(w, o, 200)
	}
}

func addTask(l Logger, f Formatter, p Parser, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t, ucerr := p.AddTask(r.Body)
		defer r.Body.Close()
		if ucerr != nil {
			l.Printf("error parsing addTask data: %v", ucerr)
			f.WriteResponse(w, f.Errorf("Error: could not parse task data: %v", ucerr), 400)
			return
		}
		td, ucerr := usecase.AddTask(taskRepo, t)
		if ucerr != nil {
			l.Printf("error adding task: %v", ucerr)
			f.WriteResponse(w, f.Error("Error: could not add task data"), 500)
			return
		}
		o, err := f.TaskID(td.TaskID)
		if err != nil {
			f.WriteResponse(w, f.Error("Task created, but there was an error formatting the response Task ID"), 201)
			return
		}
		f.WriteResponse(w, o, 201)
	}
}

func completeTask(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		taskIDInt, err := strconv.Atoi(ps.ByName("taskID"))
		if err != nil {
			l.Printf("valid task ID required")
			f.WriteResponse(w, f.Error("Error: valid task ID required"), 404)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, ucerr := usecase.CompleteTask(taskRepo, id)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Task ID %d not found", id), 404)
				return
			}
			l.Printf("error completing task: %v", ucerr)
			f.WriteResponse(w, f.Error("Error completing task"), 500)
			return
		}
		if !ok {
			f.WriteResponse(w, f.Errorf("Task %v already completed", id), 400)
			return
		}
		f.WriteEmpty(w, 204)
	}
}

func clearTask(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		taskIDInt, err := strconv.Atoi(params.ByName("taskID"))
		if err != nil {
			l.Printf("valid task ID required")
			f.WriteResponse(w, f.Error("Error: valid task ID required"), 404)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, ucerr := usecase.ClearTask(taskRepo, id)
		if ucerr != nil {
			if ucerr.Code() == usecase.ErrRecordNotFound {
				f.WriteResponse(w, f.Errorf("Task ID %d not found", id), 404)
				return
			}
			l.Printf("error clearing task: %v", ucerr)
			f.WriteResponse(w, f.Error("Error clearing task"), 500)
			return
		}
		if !ok {
			f.WriteResponse(w, f.Errorf("Task %v already cleared", id), 404)
			return
		}
		f.WriteEmpty(w, 204)
	}
}

func clearCompletedTasks(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		count, ucerr := usecase.ClearCompletedTasks(taskRepo)
		if ucerr != nil {
			l.Printf("error clearing completed tasks: %v", ucerr)
			f.WriteResponse(w, f.Error("Error clearing completed tasks"), 500)
			return
		}
		o, err := f.ClearedCompleted(count)
		if err != nil {
			f.WriteResponse(w, f.Error("Completed tasks cleared, but there was an error formatting the response"), 200)
			return
		}
		f.WriteResponse(w, o, 200)
	}
}
