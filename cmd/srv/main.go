package main

import (
	"log"
	"os"

	data "github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	"github.com/benjohns1/scheduled-tasks/internal/infra/scheduler"
	"github.com/benjohns1/scheduled-tasks/internal/present/restapi"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stderr, "main ", log.LstdFlags)

	// Load environment vars
	godotenv.Load("../../.env")

	// Load DB connection info
	dbconn := data.NewDBConn(l)
	if err := dbconn.Connect(); err != nil {
		l.Panic(err)
	}
	defer dbconn.Close()

	l.Print("starting scheduler and API server")
	scChan := startScheduler(dbconn)
	acChan := startAPIServer(dbconn)

	sc := false
	ac := false
	for {
		select {
		case sc = <-scChan:
			l.Print("scheduler closed")
		case ac = <-acChan:
			l.Print("api server closed")
		}
		if sc && ac {
			l.Print("all processes closed, exiting")
			return
		}
	}
}

func startAPIServer(dbconn data.DBConn) (closed <-chan bool) {
	l := log.New(os.Stderr, "api ", log.LstdFlags)

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

	// Serve REST API
	api := restapi.New(l, taskRepo, scheduleRepo)
	closed = restapi.Serve(l, api)
	return
}

func startScheduler(dbconn data.DBConn) (closed <-chan bool) {
	l := log.New(os.Stderr, "sched ", log.LstdFlags)

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

	// Instantiate time clock
	c := &scheduler.Clock{}

	// Start scheduler process
	_, closed, _ = scheduler.Run(l, c, taskRepo, scheduleRepo)
	return
}
