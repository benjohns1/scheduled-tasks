package main

import (
	"log"
	"os"

	data "github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	"github.com/benjohns1/scheduled-tasks/internal/infra/scheduler"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stderr, "sched ", log.LstdFlags)

	// Load environment vars
	godotenv.Load("../../.env")

	// Load DB connection info
	dbconn := data.NewDBConn(l)
	if err := dbconn.Connect(); err != nil {
		l.Panic(err)
	}
	defer dbconn.Close()
	didSetup, err := dbconn.Setup()
	if err != nil {
		l.Panic(err)
	}
	if didSetup {
		l.Print("first-time DB setup complete")
	}

	// Instantiate repositories
	taskRepo, err := data.NewTaskRepo(dbconn)
	if err != nil {
		l.Panic(err)
	}
	scheduleRepo, err := data.NewScheduleRepo(dbconn)
	if err != nil {
		l.Panic(err)
	}

	// Start scheduler process
	scheduler.Run(l, taskRepo, scheduleRepo)
}
