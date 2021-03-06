package task

import (
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
	responseMapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	mapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/task/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter defines the formatter interface for output responses
type Formatter interface {
	ClearedCompleted(count int) ([]byte, error)
	TaskID(id usecase.TaskID) ([]byte, error)
	Task(td *usecase.TaskData) ([]byte, error)
	TaskMap(ts map[usecase.TaskID]*task.Task) ([]byte, error)
	responseMapper.ResponseFormatter
}

// Parser defines the parser interface for parsing input requests
type Parser interface {
	AddTask(b io.Reader, uid user.ID) (*task.Task, error)
}

// Handle adds task handling endpoints
func Handle(r *httprouter.Router, prefix string, l Logger, rf responseMapper.ResponseFormatter, taskRepo usecase.TaskRepo) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(rf)

	pre := prefix + "/task"
	r.GET(pre+"/", auth.HRAuthorize(auth.PermReadTask, false, l, f, listTasks(l, f, taskRepo)))
	r.GET(pre+"/:taskID", auth.HRAuthorize(auth.PermReadTask, false, l, f, getTask(l, f, taskRepo)))
	r.POST(pre+"/", auth.HRAuthorize(auth.PermUpsertTask, true, l, f, addTask(l, f, p, taskRepo)))
	r.PUT(pre+"/:taskID/complete", auth.HRAuthorize(auth.PermUpsertTask, true, l, f, completeTask(l, f, taskRepo)))
	r.DELETE(pre+"/:taskID", auth.HRAuthorize(auth.PermDeleteTask, true, l, f, clearTask(l, f, taskRepo)))
	r.POST(pre+"/clear", auth.HRAuthorize(auth.PermDeleteTask, true, l, f, clearCompletedTasks(l, f, taskRepo)))
}

func listTasks(l Logger, f Formatter, taskRepo usecase.TaskRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		u := auth.GetUser(w)
		ts, ucerr := usecase.ListTasks(taskRepo, u.ID())
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
		u := auth.GetUser(w)
		id := usecase.TaskID(taskIDInt)
		td, ucerr := usecase.GetTask(taskRepo, id, u.ID())
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
		u := auth.GetUser(w)
		if u.ID().IsEmpty() {
			l.Printf("error retrieving required user from http.ResponseWriter: %v", w)
			f.ErrUnauthorized(w)
			return
		}
		t, ucerr := p.AddTask(r.Body, u.ID())
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
		u := auth.GetUser(w)
		uid := u.ID()
		if uid.IsEmpty() {
			l.Printf("error retrieving required user from http.ResponseWriter: %v", w)
			f.ErrUnauthorized(w)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, ucerr := usecase.CompleteTask(taskRepo, id, uid)
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
		u := auth.GetUser(w)
		uid := u.ID()
		if uid.IsEmpty() {
			l.Printf("error retrieving required user from http.ResponseWriter: %v", w)
			f.ErrUnauthorized(w)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, ucerr := usecase.ClearTask(taskRepo, id, uid)
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
		u := auth.GetUser(w)
		uid := u.ID()
		if uid.IsEmpty() {
			l.Printf("error retrieving required user from http.ResponseWriter: %v", w)
			f.ErrUnauthorized(w)
			return
		}
		count, ucerr := usecase.ClearCompletedTasks(taskRepo, uid)
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
