package restapi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	mapper "github.com/benjohns1/scheduled-tasks/internal/present/restapi/json"
	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Serve creates and starts the REST API server
func Serve(l Logger, taskRepo usecase.TaskRepo) {

	// Start API server
	port := 8080
	if val, err := strconv.Atoi(os.Getenv("APPLICATION_PORT")); err == nil {
		port = val
	}
	l.Printf("starting server on port %d", port)

	p := mapper.NewParser()
	f := mapper.NewFormatter(l)
	r := httprouter.New()
	taskPrefix := "/api/v1/task"
	r.GET(taskPrefix+"/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	})
	r.GET(taskPrefix+"/:taskID", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	})
	r.POST(taskPrefix+"/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	})
	r.PUT(taskPrefix+"/:taskID/complete", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	})
	r.DELETE(taskPrefix+"/:taskID", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
	})
	r.POST(taskPrefix+"/clear", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	l.Printf("server exiting")
}
