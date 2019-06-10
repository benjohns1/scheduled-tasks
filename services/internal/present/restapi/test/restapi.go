package test

import (
	"net/http"

	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	postgres_test "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi"
)

// Tester describes an API struct used to create new test API instances
type Tester interface {
	NewAPI() http.Handler
	Close() error
}

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {}

// Strp returns a pointer to the passed-in string
func Strp(str string) *string {
	return &str
}

type transientTester struct{}

func (m *transientTester) NewAPI() http.Handler {
	l := &loggerStub{}
	taskRepo := transient.NewTaskRepo()
	scheduleRepo := transient.NewScheduleRepo()
	c := make(chan<- bool)
	return restapi.New(l, c, taskRepo, scheduleRepo)
}

func (m *transientTester) Close() error {
	return nil
}

// NewTransientTester returns a tester struct for creating APIs with a transient in-memory DB
func NewTransientTester() Tester {
	return &transientTester{}
}

type postgresTester struct {
	prevConn *postgres.DBConn
}

func (m *postgresTester) NewAPI() http.Handler {
	m.Close()
	conn, err := postgres_test.NewTestDBConn(postgres_test.IntegrationTest)
	if err != nil {
		panic(err)
	}
	m.prevConn = &conn
	taskRepo, err := postgres.NewTaskRepo(conn)
	if err != nil {
		panic(err)
	}
	scheduleRepo, err := postgres.NewScheduleRepo(conn)
	if err != nil {
		panic(err)
	}
	l := &loggerStub{}
	c := make(chan<- bool)
	return restapi.New(l, c, taskRepo, scheduleRepo)
}

func (m *postgresTester) Close() error {
	if m.prevConn != nil {
		return m.prevConn.Close()
	}
	return nil
}

// NewPostgresTester returns a tester struct for creating APIs with a Postgres DB
func NewPostgresTester() Tester {
	return &postgresTester{}
}
