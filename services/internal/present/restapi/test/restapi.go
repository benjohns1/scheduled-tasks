package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	postgres_test "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Tester describes an API struct used to create new test API instances
type Tester interface {
	NewAPI() MockAPI
	Close() error
}

// MockAPI contain the API mock and repos used during setup
type MockAPI struct {
	API          http.Handler
	UserRepo     usecase.UserRepo
	TaskRepo     usecase.TaskRepo
	ScheduleRepo usecase.ScheduleRepo
}

type loggerStub struct{}

func (l *loggerStub) Printf(format string, v ...interface{}) {
	if testing.Verbose() {
		fmt.Printf(fmt.Sprintf("    LOG: %v\n", format), v...)
	}
}

// Strp returns a pointer to the passed-in string
func Strp(str string) *string {
	return &str
}

type transientTester struct{}

func (m *transientTester) NewAPI() MockAPI {
	l := &loggerStub{}
	userRepo := transient.NewUserRepo()
	taskRepo := transient.NewTaskRepo()
	scheduleRepo := transient.NewScheduleRepo()
	c := make(chan<- bool)
	authMock := NewAuthMock(l)
	api := restapi.New(l, authMock, c, userRepo, taskRepo, scheduleRepo)
	return MockAPI{api, userRepo, taskRepo, scheduleRepo}
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

func (m *postgresTester) NewAPI() MockAPI {
	m.Close()
	conn, err := postgres_test.NewTestDBConn(postgres_test.IntegrationTest)
	if err != nil {
		panic(err)
	}
	m.prevConn = &conn

	userRepo, err := postgres.NewUserRepo(conn)
	if err != nil {
		panic(err)
	}
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
	authMock := NewAuthMock(l)
	api := restapi.New(l, authMock, c, userRepo, taskRepo, scheduleRepo)
	return MockAPI{api, userRepo, taskRepo, scheduleRepo}
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
