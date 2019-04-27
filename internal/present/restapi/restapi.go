package restapi

import (
	"fmt"
	"log"
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
	log.Printf("starting server on port %d", port)

	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ts, err := usecase.ListTasks(taskRepo)
		if err != nil {
			log.Printf("error retrieving tasks")
			w.Write([]byte("Error: couldn't retrieve tasks"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Found %d tasks: %v", len(ts), ts)))
	})
	r.GET("/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		task, err := usecase.AddTask(taskRepo, "name", "desc")
		if err != nil {
			log.Printf("error adding task: %v", err)
			w.Write([]byte("Error adding task"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Added task %v: %v", task.TaskID, task)))
	})
	r.GET("/complete/:taskID", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		taskIDInt, err := strconv.Atoi(ps.ByName("taskID"))
		if err != nil {
			log.Printf("valid task id required")
			w.Write([]byte("Error: valid task ID required"))
			return
		}

		id := usecase.TaskID(taskIDInt)
		ok, err := usecase.CompleteTask(taskRepo, id)
		if err != nil {
			log.Printf("error completing task: %v", err)
			w.Write([]byte("Error completing task"))
			return
		}
		if ok {
			w.Write([]byte(fmt.Sprintf("Completed task %v", id)))
			return
		}
		w.Write([]byte(fmt.Sprintf("Task %v already completed", id)))
	})
	r.GET("/clear", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		count, err := usecase.ClearCompletedTasks(taskRepo)
		if err != nil {
			log.Printf("error clearing completed tasks: %v", err)
			w.Write([]byte("Error clearing completed tasks"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Cleared %d completed tasks", count)))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	log.Printf("server exiting")
}
