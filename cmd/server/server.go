package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/benjohns1/scheduled-tasks/internal/app/taskapp"
	"github.com/benjohns1/scheduled-tasks/internal/pkg/persistence"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi"
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

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world"))
	})
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		task := &taskapp.Task{Name: "asdf", Description: "description"}
		key, err := taskapp.AddTask(db, task)
		if err != nil {
			log.Printf("error adding task: %v", err)
			w.Write([]byte("Error adding task"))
		}
		w.Write([]byte(fmt.Sprintf("Added task %v: %v", key, task)))
	})
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	log.Printf("server exiting")
}
