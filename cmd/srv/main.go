package main

import (
	"log"
	"os"

	data "github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	"github.com/benjohns1/scheduled-tasks/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/internal/present/restapi"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stderr, "api ", log.LstdFlags)

	// Load environment vars
	godotenv.Load("../../.env")

	// Load DB connection info
	dbconn := data.NewDBConn()
	taskRepo, err := data.NewTaskRepo(l, dbconn)
	scheduleRepo := transient.NewScheduleRepo()
	if err != nil {
		l.Panic(err)
	}
	defer taskRepo.Close()

	restapi.Serve(l, taskRepo, scheduleRepo)
}
