// +build integration

package restapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/schedule"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/task"
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/test"
)

func TestTransientRESTAPIBasic(t *testing.T) {
	tester := test.NewTransientTester()
	defer tester.Close()
	suiteBasic(t, tester)
}

func TestPostgresRESTAPIBasic(t *testing.T) {
	tester := test.NewPostgresTester()
	defer tester.Close()
	suiteBasic(t, tester)
}

func suiteBasic(t *testing.T, tester test.Tester) {
	addOrUpdateExternalUser(t, tester.NewAPI())
	errorResponse(t, tester.NewAPI())
	listTasks(t, tester.NewAPI())
	addTask(t, tester.NewAPI())
	getTask(t, tester.NewAPI())
	completeTask(t, tester.NewAPI())
	clearTask(t, tester.NewAPI())
	clearCompletedTasks(t, tester.NewAPI())
	listSchedules(t, tester.NewAPI())
	addRecurringTask(t, tester.NewAPI())
	addSchedule(t, tester.NewAPI())
	getSchedule(t, tester.NewAPI())
	pauseSchedule(t, tester.NewAPI())
	unpauseSchedule(t, tester.NewAPI())
	removeSchedule(t, tester.NewAPI())
}

func addOrUpdateExternalUser(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API
	u1 := user.New("test user")
	apiMock.UserRepo.AddExternal(u1, "https://p1/", "e1")
	u1Perms := []auth.Permission{auth.PermUpsertUserSelf}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "https://p1/", Subject: "e1", Permissions: u1Perms}, api)

	u2 := user.New("test user, no perms")
	apiMock.UserRepo.AddExternal(u2, "https://p1/", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "https://p1/", Subject: "e2"}, api)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "should return 204 for a valid user and token",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/user/external/p1/e1/addOrUpdate", body: `{"displayname":"update my username"}`},
			asserts: asserts{statusEquals: http.StatusNoContent, bodyEquals: test.Strp(``)},
		},
		{
			name:    "should return 401 for valid user without permissions",
			h:       u2Api,
			args:    args{method: "PUT", url: "/api/v1/user/external/p1/e2/addOrUpdate", body: `{"displayname":"badPerms!"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "should return 401 for invalid user",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/user/external/p1/invalid-user-id/addOrUpdate", body: `{"displayname":"myNameIsWho?"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "should return 401 for unauthorized user",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/user/external/someProvider/userExternalID/addOrUpdate", body: `{"displayname":"myNameIsWho?"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func errorResponse(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "should return 404 with JSON error",
			h:       api,
			args:    args{method: "GET", url: "/invalid-resource-uri"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyEquals: test.Strp(`{"error":"Not found"}`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func listTasks(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	nowStr, reset := test.SetStaticClock(time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC))
	defer reset()

	u1 := user.New("test user")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermReadTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)

	u2 := user.New("test user, no perms")
	apiMock.UserRepo.AddExternal(u2, "p1", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e2"}, api)

	u3 := user.New("test user with tasks")
	apiMock.UserRepo.AddExternal(u3, "p1", "e3")
	u3t1 := task.New("u3 task1", "u3t1 description", u3.ID())
	apiMock.TaskRepo.Add(u3t1)
	u3Perms := []auth.Permission{auth.PermReadTask}
	u3Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e3", Permissions: u3Perms}, api)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "u1 should return 200 empty list",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
		},
		{
			name:    "u2 invalid permissions should return 401",
			h:       u2Api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "u3 should return list with 1 task",
			h:       u3Api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(fmt.Sprintf(`{"1":{"id":1,"name":"u3 task1","description":"u3t1 description","completedTime":null,"createdTime":"%v"}}`, nowStr))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func addTask(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1 := user.New("test user for addTask")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermUpsertTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)

	u2 := user.New("test user for addTask, no perms")
	apiMock.UserRepo.AddExternal(u2, "p1", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e2"}, api)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "empty task should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "invalid JSON should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{{{`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse task data`)},
		},
		{
			name:    "empty body should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/task/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse task data`)},
		},
		{
			name:    "invalid permissions should return 401",
			h:       u2Api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}
func getTask(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	nowStr, reset := test.SetStaticClock(time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC))
	defer reset()

	u1 := user.New("test user for getTask")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermReadTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)
	u1t1 := task.New("u1 task1", "u1t1 task description", u1.ID())
	apiMock.TaskRepo.Add(u1t1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown task ID should return 404",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/task/9999"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 9999 not found`)},
		},
		{
			name:    "should return valid task",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(fmt.Sprintf(`{"id":1,"name":"u1 task1","description":"u1t1 task description","completedTime":null,"createdTime":"%v"}`, nowStr))},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func completeTask(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1 := user.New("test user for completeTask")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermUpsertTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)
	u1t1 := task.New("u1t1 task", "", u1.ID())
	apiMock.TaskRepo.Add(u1t1)

	u2 := user.New("test user for completeTask, no perms")
	apiMock.UserRepo.AddExternal(u2, "p1", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e2"}, api)
	u2t1 := task.New("u2t1 task", "", u2.ID())
	apiMock.TaskRepo.Add(u2t1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/task/1/complete"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown task ID should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/task/9999/complete"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 9999 not found`)},
		},
		{
			name:    "valid task ID should return 204",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/task/1/complete"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "valid task ID without proper permissions should return 401",
			h:       u2Api,
			args:    args{method: "PUT", url: "/api/v1/task/2/complete"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid task ID owned by another user should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/task/2/complete"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 2 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func clearTask(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1 := user.New("test user for clearTask")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermDeleteTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)
	u1t1 := task.New("u1t1 task", "", u1.ID())
	apiMock.TaskRepo.Add(u1t1)

	u2 := user.New("test user for clearTask, no perms")
	apiMock.UserRepo.AddExternal(u2, "p1", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e2"}, api)
	u2t1 := task.New("u2t1 task", "", u2.ID())
	apiMock.TaskRepo.Add(u2t1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "DELETE", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown task ID should return 404",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/task/9999"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 9999 not found`)},
		},
		{
			name:    "valid task ID should return 204",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "valid task ID without proper permissions should return 401",
			h:       u2Api,
			args:    args{method: "DELETE", url: "/api/v1/task/2"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid task ID owned by another user should return 404",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/task/2"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 2 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func clearCompletedTasks(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1 := user.New("test user 1 for clearTask")
	apiMock.UserRepo.AddExternal(u1, "p1", "e1")
	u1Perms := []auth.Permission{auth.PermDeleteTask}
	u1Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e1", Permissions: u1Perms}, api)
	u1t1 := task.New("u1t1 task", "", u1.ID())
	u1t1.CompleteNow()
	u1t2 := task.New("u1t2 task", "", u1.ID())
	u1t2.CompleteNow()
	apiMock.TaskRepo.Add(u1t1)
	apiMock.TaskRepo.Add(u1t2)

	u2 := user.New("test user 2 for clearTask, no perms")
	apiMock.UserRepo.AddExternal(u2, "p1", "e2")
	u2Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e2"}, api)
	u2t1 := task.New("u2t1 task", "", u2.ID())
	u2t1.CompleteNow()
	apiMock.TaskRepo.Add(u2t1)

	u3 := user.New("test user 3 for clearTask")
	apiMock.UserRepo.AddExternal(u3, "p1", "e3")
	u3Perms := []auth.Permission{auth.PermDeleteTask}
	u3Api := test.InjectClaims(test.MockClaims{Issuer: "p1", Subject: "e3", Permissions: u3Perms}, api)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/clear"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "clearing task list should return 200 with 2 count",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/task/clear"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"count":2,"message":"Cleared all completed tasks"}`)},
		},
		{
			name:    "clearing task list with invalid permissions should return 401",
			h:       u2Api,
			args:    args{method: "POST", url: "/api/v1/task/clear"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "clearing empty task list should return 200 with 0 count",
			h:       u3Api,
			args:    args{method: "POST", url: "/api/v1/task/clear"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"count":0,"message":"No completed tasks to clear"}`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func listSchedules(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	_, u1Api := apiMock.NewUserWithPerm("test user for listSchedules", "p1", "e1", auth.PermReadSchedule)

	u2, u2Api := apiMock.NewUserWithPerm("test user for listSchedules, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	u3, u3Api := apiMock.NewUserWithPerm("test user for listSchedules, with tasks", "p1", "e3", auth.PermReadSchedule)
	u3f1, _ := schedule.NewHourFrequency([]int{0, 30})
	u3s1 := schedule.New(u3f1, u3.ID())
	apiMock.ScheduleRepo.Add(u3s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "u1 should return 200 empty list",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
		},
		{
			name:    "u2 invalid permissions should return 401",
			h:       u2Api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "u3 should return list with 1 schedule",
			h:       u3Api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"2":{"id":2,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,30],"paused":false,"tasks":[]}}`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func addSchedule(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	_, u1Api := apiMock.NewUserWithPerm("test user for addSchedule", "p1", "e1", auth.PermUpsertSchedule)
	_, u2Api := apiMock.NewUserWithPerm("test user for addSchedule, no perms", "p1", "e2", auth.PermNone)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "empty hour schedule should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "empty day schedule should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Day"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "empty week schedule should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Week"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":3}`)},
		},
		{
			name:    "empty month schedule should return 201 and ID",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Month"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":4}`)},
		},
		{
			name:    "empty/invalid schedule should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data: invalid frequency`)},
		},
		{
			name:    "invalid JSON should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{{{`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data`)},
		},
		{
			name:    "empty body should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data`)},
		},
		{
			name:    "invalid permissions should return 401",
			h:       u2Api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func getSchedule(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1, u1Api := apiMock.NewUserWithPerm("test user for getSchedule", "p1", "e1", auth.PermReadSchedule)
	u1f1, _ := schedule.NewHourFrequency([]int{0})
	u1s1 := schedule.New(u1f1, u1.ID())
	apiMock.ScheduleRepo.Add(u1s1)

	u2, u2Api := apiMock.NewUserWithPerm("test user for getSchedule, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid ID should return schedule data",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "other user's schedule should return 404",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusNotFound},
		},
		{
			name:    "user without permissions should return 401",
			h:       u2Api,
			args:    args{method: "GET", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown ID should return 404",
			h:       u1Api,
			args:    args{method: "GET", url: "/api/v1/schedule/9999"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 9999 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func removeSchedule(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1, u1Api := apiMock.NewUserWithPerm("test user for removeSchedule", "p1", "e1", auth.PermDeleteSchedule)
	u1f1, _ := schedule.NewHourFrequency([]int{0})
	u1s1 := schedule.New(u1f1, u1.ID())
	apiMock.ScheduleRepo.Add(u1s1)

	u2, u2Api := apiMock.NewUserWithPerm("test user for removeSchedule, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid ID should return 204",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "other user's schedule should return 404",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 2 not found`)},
		},
		{
			name:    "user without permissions should return 401",
			h:       u2Api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown ID should return 404",
			h:       u1Api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/9999"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 9999 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func pauseSchedule(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1, u1Api := apiMock.NewUserWithPerm("test user for removeSchedule", "p1", "e1", auth.PermUpsertSchedule)
	u1f1, _ := schedule.NewHourFrequency([]int{0})
	u1s1 := schedule.New(u1f1, u1.ID())
	apiMock.ScheduleRepo.Add(u1s1)

	u2, u2Api := apiMock.NewUserWithPerm("test user for removeSchedule, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/pause"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid ID should return 204",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/pause"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "other user's schedule should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/2/pause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 2 not found`)},
		},
		{
			name:    "user without permissions should return 401",
			h:       u2Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/2/pause"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown ID should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/9999/pause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 9999 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func unpauseSchedule(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1, u1Api := apiMock.NewUserWithPerm("test user for removeSchedule", "p1", "e1", auth.PermUpsertSchedule)
	u1f1, _ := schedule.NewHourFrequency([]int{0})
	u1s1 := schedule.New(u1f1, u1.ID())
	apiMock.ScheduleRepo.Add(u1s1)

	u2, u2Api := apiMock.NewUserWithPerm("test user for removeSchedule, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/unpause"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid ID should return 204",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/unpause"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "other user's schedule should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/2/unpause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 2 not found`)},
		},
		{
			name:    "user without permissions should return 401",
			h:       u2Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/2/unpause"},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown ID should return 404",
			h:       u1Api,
			args:    args{method: "PUT", url: "/api/v1/schedule/9999/unpause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 9999 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}

func addRecurringTask(t *testing.T, apiMock test.MockAPI) {
	api := apiMock.API

	u1, u1Api := apiMock.NewUserWithPerm("test user for removeSchedule", "p1", "e1", auth.PermUpsertSchedule)
	u1f1, _ := schedule.NewHourFrequency([]int{0})
	u1s1 := schedule.New(u1f1, u1.ID())
	apiMock.ScheduleRepo.Add(u1s1)

	u2, u2Api := apiMock.NewUserWithPerm("test user for removeSchedule, no perms", "p1", "e2", auth.PermNone)
	u2f1, _ := schedule.NewHourFrequency([]int{0})
	u2s1 := schedule.New(u2f1, u2.ID())
	apiMock.ScheduleRepo.Add(u2s1)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "no auth should return 401",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{"name": "t1", "description": "t1 desc"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "valid schedule ID should return 201",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{"name": "t1", "description": "t1 desc"}`},
			asserts: asserts{statusEquals: http.StatusCreated},
		},
		{
			name:    "valid schedule ID with empty task body should return 201",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusCreated},
		},
		{
			name:    "valid schedule ID with null body should return 400",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest},
		},
		{
			name:    "other user's schedule should return 404",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/2/task/", body: `{"name": "t1", "description": "t1 desc"}`},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 2 not found`)},
		},
		{
			name:    "user without permissions should return 401",
			h:       u2Api,
			args:    args{method: "POST", url: "/api/v1/schedule/2/task/", body: `{"name": "t1", "description": "t1 desc"}`},
			asserts: asserts{statusEquals: http.StatusUnauthorized},
		},
		{
			name:    "unknown ID should return 404",
			h:       u1Api,
			args:    args{method: "POST", url: "/api/v1/schedule/9999/task/", body: `{"name": "t1", "description": "t1 desc"}`},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 9999 not found`)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if tt.asserts.bodyEquals != nil && rr.Body.String() != *tt.asserts.bodyEquals {
				t.Errorf("response body = %v, should equal %v", rr.Body.String(), *tt.asserts.bodyEquals)
			}
			if tt.asserts.bodyContains != nil && !strings.Contains(rr.Body.String(), *tt.asserts.bodyContains) {
				t.Errorf("response body = %v, should contain %v", rr.Body.String(), *tt.asserts.bodyContains)
			}
		})
	}
}
