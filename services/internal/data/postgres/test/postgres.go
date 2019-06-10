package test

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	"github.com/joho/godotenv"
)

type loggerStub struct{}

func (l *loggerStub) Print(v ...interface{})                 {}
func (l *loggerStub) Printf(format string, v ...interface{}) {}
func (l *loggerStub) Println(v ...interface{})               {}

type testType uint8

// Test type
const (
	DBTest testType = iota
	IntegrationTest
)

// NewTestDBConn creates a fake Postgres DB connection
func NewTestDBConn(test testType) (postgres.DBConn, error) {
	l := &loggerStub{}

	// Load environment vars
	if err := godotenv.Load("../../../../.env"); err != nil {
		return postgres.DBConn{}, fmt.Errorf("could not load .env file: %v", err)
	}

	var portEnvVar string
	switch test {
	case DBTest:
		portEnvVar = "POSTGRES_DBTEST_PORT"
	case IntegrationTest:
		portEnvVar = "POSTGRES_INTEGRATION_PORT"
	}

	// Load DB connection info
	dbconn := postgres.NewDBConn(l)
	testPort, err := strconv.Atoi(os.Getenv(portEnvVar))
	if err != nil || testPort == 0 {
		dbconn.Close()
		return dbconn, fmt.Errorf("POSTGRES_TEST_PORT must be set to run postgres DB integreation tests %v", err)
	}
	dbconn.Port = testPort
	dbconn.MaxRetryAttempts = 1
	dbconn.RetrySleepSeconds = 0

	// Connect to DB, destroy any existing tables, setup tables again
	if err := dbconn.Connect(); err != nil {
		dbconn.Close()
		return dbconn, fmt.Errorf("could not connect to test DB: %v", err)
	}
	_, destroyErr := destroy(&dbconn)
	didSetup, err := dbconn.Setup()
	if err != nil {
		dbconn.Close()
		return dbconn, fmt.Errorf("error setting up DB tables: %v", err)
	}
	if !didSetup {
		dbconn.Close()
		return dbconn, fmt.Errorf("could not setup fresh DB tables, test tables may not have been not properly destroyed: %v", destroyErr)
	}
	return dbconn, nil
}

// destroy !!!WARNING!!! completely destroys all data in the DB
func destroy(conn *postgres.DBConn) (sql.Result, error) {
	return conn.DB.Exec("DROP TABLE task; DROP TABLE recurring_task; DROP TABLE schedule;")
}
