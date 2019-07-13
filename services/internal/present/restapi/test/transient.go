package test

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi"
)

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
