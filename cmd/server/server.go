package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benjohns1/scheduled-tasks/internal/app/taskapp"
	persistence "github.com/benjohns1/scheduled-tasks/internal/pkg/persistence/postgres"
	"github.com/joho/godotenv"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment vars
	godotenv.Load("../../.env")

	// Load DB connection info
	connInfo, err := persistence.LoadConnInfo()
	if err != nil {
		log.Panicf("error loading db connection details: %v", err)
	}

	// Connect to DB
	log.Printf("connecting to db %s as %s...", connInfo.Name, connInfo.User)
	db, err := persistence.Connect(connInfo)
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		log.Panicf("error opening db: %v", err)
	}

	// Perform DB setup if needed
	setup, err := persistence.Setup(db)
	if err != nil {
		log.Panicf("error setting up db: %v", err)
	}
	if setup {
		log.Print("first-time DB setup complete")
	}

	// Start API server
	port := 8080
	if val, err := strconv.Atoi(os.Getenv("APPLICATION_PORT")); err == nil {
		port = val
	}
	log.Printf("starting server on port %d", port)

	r := httprouter.New()
	r.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("Hello, world"))
	})
	r.GET("/add", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		task := &taskapp.Task{Name: "task name", Description: "description"}
		id, err := taskapp.AddTask(db, task)
		if err != nil {
			log.Printf("error adding task: %v", err)
			w.Write([]byte("Error adding task"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Added task %v: %v", id, task)))
	})
	r.GET("/complete/:taskID", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		taskIDInt, err := strconv.Atoi(ps.ByName("taskID"))
		if err != nil {
			log.Printf("valid task id required")
			w.Write([]byte("Error: valid task ID required"))
			return
		}

		id := taskapp.TaskID(taskIDInt)
		err = taskapp.CompleteTask(db, id)
		if err != nil {
			log.Printf("error completing task: %v", err)
			w.Write([]byte("Error completing task"))
			return
		}
		w.Write([]byte(fmt.Sprintf("Completed task %v", id)))
	})
	r.GET("/clear", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		count, err := taskapp.ClearCompleted(db)
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