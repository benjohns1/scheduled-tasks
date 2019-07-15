package main

import (
	"log"
	"os"

	data "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	"github.com/benjohns1/scheduled-tasks/services/internal/infra/scheduler"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
	"github.com/joho/godotenv"
)

func main() {
	l := log.New(os.Stderr, "main ", log.LstdFlags)

	loadEnv()

	// Scheduler DB connection
	scConn := data.NewDBConn(l, "scheduler")
	if err := scConn.Connect(); err != nil {
		l.Panic(err)
	}
	defer scConn.Close()

	// API DB connection
	acConn := data.NewDBConn(l, "api")
	if err := acConn.Connect(); err != nil {
		l.Panic(err)
	}
	defer acConn.Close()

	l.Print("starting scheduler and API server")
	checkC, scChan := startScheduler(scConn)
	acChan := startAPIServer(acConn, checkC)

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

func loadEnv() {
	// Load environment vars from file, if ENV_FILEPATH was set
	if path, exists := os.LookupEnv("ENV_FILEPATH"); exists {
		godotenv.Load(path)
	}
}

func startAPIServer(dbconn data.DBConn, check chan<- bool) (closed <-chan bool) {
	l := log.New(os.Stderr, "api ", log.LstdFlags)

	didSetup, err := dbconn.Setup()
	if err != nil {
		l.Panic(err)
	}
	if didSetup {
		l.Print("first-time DB setup complete")
	}

	// Instantiate repositories
	userRepo, err := data.NewUserRepo(dbconn)
	if err != nil {
		l.Panic(err)
	}
	taskRepo, err := data.NewTaskRepo(dbconn)
	if err != nil {
		l.Panic(err)
	}
	scheduleRepo, err := data.NewScheduleRepo(dbconn)
	if err != nil {
		l.Panic(err)
	}

	// Instantiate authorization handler
	a := auth.NewAuth0(l, auth.Auth0Config{
		Secret:   []byte(os.Getenv("AUTH0_API_SECRET")),
		Audience: []string{os.Getenv("AUTH0_API_IDENTIFIER")},
		Domain:   os.Getenv("AUTH0_DOMAIN"),
	})

	// Serve REST API
	api := restapi.New(l, a, check, userRepo, taskRepo, scheduleRepo)
	return restapi.Serve(l, api)
}

func startScheduler(dbconn data.DBConn) (check chan<- bool, closed <-chan bool) {
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

	// Start scheduler process
	_, check, closed = scheduler.Run(l, taskRepo, scheduleRepo, nil)
	return check, closed
}
