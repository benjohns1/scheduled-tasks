package main

import (
	"log"

	data "github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	"github.com/benjohns1/scheduled-tasks/internal/present/restapi"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment vars
	godotenv.Load("../../.env")

	// Load DB connection info
	dbconn := data.NewDBConn()
	taskRepo, err := data.NewTaskRepo(dbconn)
	if err != nil {
		log.Panic(err)
	}
	defer taskRepo.Close()

	restapi.Serve(taskRepo)
}
