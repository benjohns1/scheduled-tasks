package restapi

import (
	"fmt"
	logger "log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/internal/usecase"
)

// Serve creates and starts the REST API server
func Serve(taskRepo usecase.TaskRepo) {

	// Start API server
	port := 8080
	if val, err := strconv.Atoi(os.Getenv("APPLICATION_PORT")); err == nil {
		port = val
	}
	log := logger.New(os.Stderr, "api ", logger.LstdFlags)
	log.Printf("starting server on port %d", port)

	r := httprouter.New()
	taskPrefix := "/api/v1/task"
	r.GET(taskPrefix+"/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ts, err := usecase.ListTasks(taskRepo)
		if err != nil {
			log.Printf("error retrieving task list: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error: couldn't retrieve tasks")), 500)
			return
		}
		o, err := taskMapToJSON(ts)
		if err != nil {
			log.Printf("error encoding task map: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error encoding task data")), 500)
		}
		writeResponse(w, o, 200)
	})
	r.POST(taskPrefix+"/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		t, err := addTaskFromJSON(r.Body)
		if err != nil {
			log.Printf("error parsing addTask data: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error: could not parse task data: %v", err)), 400)
			return
		}
		td, err := usecase.AddTask(taskRepo, t)
		if err != nil {
			log.Printf("error adding task: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error: could not add task data")), 500)
			return
		}
		o, err := idToJSON(td.TaskID)
		if err != nil {
			writeResponse(w, errorToJSON(log, fmt.Errorf("Task created, but there was an error formatting the response Task ID")), 201)
			return
		}
		writeResponse(w, o, 201)
	})
	r.PUT(taskPrefix+"/:taskID/complete", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		taskIDInt, err := strconv.Atoi(ps.ByName("taskID"))
		if err != nil {
			log.Printf("valid task ID required")
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error: valid task ID required")), 404)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, err := usecase.CompleteTask(taskRepo, id)
		if err != nil {
			log.Printf("error completing task: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error completing task")), 500)
			return
		}
		if !ok {
			writeResponse(w, errorToJSON(log, fmt.Errorf("Task %v already completed", id)), 400)
			return
		}
		writeEmpty(w, 204)
	})
	r.DELETE(taskPrefix+"/:taskID", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		taskIDInt, err := strconv.Atoi(ps.ByName("taskID"))
		if err != nil {
			log.Printf("valid task ID required")
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error: valid task ID required")), 404)
			return
		}
		id := usecase.TaskID(taskIDInt)
		ok, err := usecase.ClearTask(taskRepo, id)
		if err != nil {
			log.Printf("error clearing task: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error clearing task")), 500)
			return
		}
		if !ok {
			writeResponse(w, errorToJSON(log, fmt.Errorf("Task %v already cleared", id)), 404)
			return
		}
		writeEmpty(w, 204)
	})
	r.POST(taskPrefix+"/clear", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		count, err := usecase.ClearCompletedTasks(taskRepo)
		if err != nil {
			log.Printf("error clearing completed tasks: %v", err)
			writeResponse(w, errorToJSON(log, fmt.Errorf("Error clearing completed taskS")), 500)
			return
		}
		o, err := clearedCompletedToJSON(count)
		if err != nil {
			writeResponse(w, errorToJSON(log, fmt.Errorf("Completed tasks cleared, but there was an error formatting the response")), 200)
			return
		}
		writeResponse(w, o, 200)
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	log.Printf("server exiting")
}
