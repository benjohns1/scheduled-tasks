package test

import (
	"net/http"

	"github.com/benjohns1/scheduled-tasks/internal/data/postgres"
	postgres_test "github.com/benjohns1/scheduled-tasks/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/internal/present/restapi"
)

// Mock describes the mock API struct used to create new API instances
type Mock interface {
	NewAPI() http.Handler
	Close() error
}

type loggerMock struct{}

func (l *loggerMock) Printf(format string, v ...interface{}) {}

// Strp returns a pointer to the passed-in string
func Strp(str string) *string {
	return &str
}

type mockTransient struct{}

func (m *mockTransient) NewAPI() http.Handler {
	taskRepo := transient.NewTaskRepo()
	scheduleRepo := transient.NewScheduleRepo()
	l := &loggerMock{}
	return restapi.New(l, taskRepo, scheduleRepo)
}

func (m *mockTransient) Close() error {
	return nil
}

// NewMockTransientAPI returns a new API with a transient in-memory DB
func NewMockTransientAPI() Mock {
	return &mockTransient{}
}

type mockPostgres struct {
	prevConn *postgres.DBConn
}

func (m *mockPostgres) NewAPI() http.Handler {
	m.Close()
	conn, err := postgres_test.NewMockDBConn(postgres_test.IntegrationTest)
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
	l := &loggerMock{}
	return restapi.New(l, taskRepo, scheduleRepo)
}

func (m *mockPostgres) Close() error {
	if m.prevConn != nil {
		return m.prevConn.Close()
	}
	return nil
}

// NewMockPostgresAPI returns a new API with a Postgres DB
func NewMockPostgresAPI() Mock {
	return &mockPostgres{}
}
