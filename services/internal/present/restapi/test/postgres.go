package test

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	pgtest "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi"
)

type postgresTester struct {
	prevConn *postgres.DBConn
}

func (m *postgresTester) NewAPI() MockAPI {
	m.Close()
	conn, err := pgtest.NewTestDBConn(pgtest.IntegrationTest)
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
