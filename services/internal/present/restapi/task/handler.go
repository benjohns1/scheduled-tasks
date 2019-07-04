package task

import (
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
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
	AddTask(b io.Reader) (*task.Task, error)
}

// Handle adds task handling endpoints
func Handle(r *httprouter.Router, prefix string, l Logger, rf responseMapper.ResponseFormatter, taskRepo usecase.TaskRepo) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(rf)

	tPre := prefix + "/task"
	r.GET(tPre+"/", auth.Handler(listTasks(l, f, taskRepo)))
	r.GET(tPre+"/:taskID", auth.Handler(getTask(l, f, taskRepo)))
	r.POST(tPre+"/", auth.Handler(addTask(l, f, p, taskRepo)))
	r.PUT(tPre+"/:taskID/complete", auth.Handler(completeTask(l, f, taskRepo)))
	r.DELETE(tPre+"/:taskID", auth.Handler(clearTask(l, f, taskRepo)))
	r.POST(tPre+"/clear", auth.Handler(clearCompletedTasks(l, f, taskRepo)))
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
