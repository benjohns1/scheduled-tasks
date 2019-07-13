package test

import (
	"net/http"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
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

// NewUserWithPerm creates and adds a new user and injects a mock permission claim for them in the returned http.Handler
func (m *MockAPI) NewUserWithPerm(displayname string, provider string, externalID string, perm auth.Permission) (*user.User, http.Handler) {
	return m.NewUserWithPerms(displayname, provider, externalID, []auth.Permission{perm})
}

// NewUserWithPerms creates and adds a new user and injects mock permission claims for them in the returned http.Handler
func (m *MockAPI) NewUserWithPerms(displayname string, provider string, externalID string, perms []auth.Permission) (*user.User, http.Handler) {
	u := user.New(displayname)
	m.UserRepo.AddExternal(u, provider, externalID)
	api := InjectClaims(MockClaims{Issuer: provider, Subject: externalID, Permissions: perms}, m.API)
	return u, api
}

// Strp returns a pointer to the passed-in string
func Strp(str string) *string {
	return &str
}
