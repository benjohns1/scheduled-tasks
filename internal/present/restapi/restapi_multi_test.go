// +build integration

package restapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	
	"github.com/benjohns1/scheduled-tasks/internal/present/restapi/test"
)

func TestTransientRestAPIMulti(t *testing.T) {
	m := test.NewMockTransientAPI()
	defer m.Close()
	suiteMulti(t, m)
}

func TestPostgresRestAPIMulti(t *testing.T) {
	m := test.NewMockPostgresAPI()
	defer m.Close()
	suiteMulti(t, m)
}

func suiteMulti(t *testing.T, m test.Mock) {
	addListGetCompleteTasks(t, m.NewAPI())
	addListGetScheduleAndAddRemoveRecurringTasks(t, m.NewAPI())
}

func addListGetCompleteTasks(t *testing.T, api http.Handler) {
	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains   *string
		bodyNotEquals   *string
		bodyNotContains   *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "should return 200 empty list",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task2", "description": "task2 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task3", "description": "task3 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":3}`)},
		},
		{
			name:    "should return 200 list with 3 tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"1":{"id":1,"name":"task1","description":"task1 description","completedTime":null},"2":{"id":2,"name":"task2","description":"task2 description","completedTime":null},"3":{"id":3,"name":"task3","description":"task3 description","completedTime":null}}`)},
		},
		{
			name:    "get task ID 1 should return incompleted task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"name":"task1","description":"task1 description","completedTime":null}`)},
		},
		{
			name:    "complete ID 1 should return 204",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/task/1/complete"},
			asserts: asserts{statusEquals: http.StatusNoContent, bodyContains: test.Strp(``)},
		},
		{
			name:    "get task ID 1 should return completed task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyContains: test.Strp(`{"id":1,"name":"task1","description":"task1 description","completedTime":`), bodyNotContains: test.Strp(`"completedTime":null`)},
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
			if tt.asserts.bodyNotEquals != nil && rr.Body.String() == *tt.asserts.bodyNotEquals {
				t.Errorf("response body = %v, should not equal %v", rr.Body.String(), *tt.asserts.bodyNotEquals)
			}
			if tt.asserts.bodyNotContains != nil && strings.Contains(rr.Body.String(), *tt.asserts.bodyNotContains) {
				t.Errorf("response body = %v, should not contain %v", rr.Body.String(), *tt.asserts.bodyNotContains)
			}
		})
	}
}

func addListGetScheduleAndAddRemoveRecurringTasks(t *testing.T, api http.Handler) {
	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains   *string
		bodyNotEquals   *string
		bodyNotContains   *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "should return 200 empty list",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
		},
		{
			name:    "empty/invalid schedule should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data: invalid frequency`)},
		},
		{
			name:    "empty hourly schedule should return 201 and ID 1",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"hourly", "atMinutes":[]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "new hourly schedule should return 201 and ID 2",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "hourly", "atMinutes": [0,30], "paused":true}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "new hourly schedule should return 201 and ID 3",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "hourly", "atMinutes": [0,30,59], "tasks": [{"name": "rtask1", "description": "rtask1 desc"}]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":3}`)},
		},
		{
			name:    "get schedule ID 1 should return empty schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"hourly","atMinutes":[],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 2 should return hourly schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":2,"frequency":"hourly","atMinutes":[0,30],"paused":true,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 3 should return hourly schedule with 1 recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/3"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":3,"frequency":"hourly","atMinutes":[0,30,59],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]}`)},
		},
		{
			name:    "should return 200 list with 3 schedules",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"1":{"id":1,"frequency":"hourly","atMinutes":[],"paused":false,"tasks":[]},"2":{"id":2,"frequency":"hourly","atMinutes":[0,30],"paused":true,"tasks":[]},"3":{"id":3,"frequency":"hourly","atMinutes":[0,30,59],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]}}`)},
		},
		{
			name:    "adding recurring task to schedule ID 1 should return 201",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(``)},
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
			if tt.asserts.bodyNotEquals != nil && rr.Body.String() == *tt.asserts.bodyNotEquals {
				t.Errorf("response body = %v, should not equal %v", rr.Body.String(), *tt.asserts.bodyNotEquals)
			}
			if tt.asserts.bodyNotContains != nil && strings.Contains(rr.Body.String(), *tt.asserts.bodyNotContains) {
				t.Errorf("response body = %v, should not contain %v", rr.Body.String(), *tt.asserts.bodyNotContains)
			}
		})
	}
}