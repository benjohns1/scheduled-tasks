// +build integration

package restapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/benjohns1/scheduled-tasks/internal/present/restapi/test"
)

func TestTransientRestAPIBasic(t *testing.T) {
	m := test.NewMockTransientAPI()
	defer m.Close()
	suiteBasic(t, m)
}

func TestPostgresRestAPIBasic(t *testing.T) {
	m := test.NewMockPostgresAPI()
	defer m.Close()
	suiteBasic(t, m)
}

func suiteBasic(t *testing.T, m test.Mock) {
	listTasks(t, m.NewAPI())
	addTask(t, m.NewAPI())
	getTask(t, m.NewAPI())
	completeTask(t, m.NewAPI())
	clearTask(t, m.NewAPI())
	clearCompletedTasks(t, m.NewAPI())
	listSchedules(t, m.NewAPI())
	addRecurringTask(t, m.NewAPI())
	addSchedule(t, m.NewAPI())
	getSchedule(t, m.NewAPI())
	pauseSchedule(t, m.NewAPI())
	unpauseSchedule(t, m.NewAPI())
}

func listTasks(t *testing.T, api http.Handler) {

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
			name:    "should return 200 empty list",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
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
func addTask(t *testing.T, api http.Handler) {

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
			name:    "empty task should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "invalid JSON should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{{{`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse task data`)},
		},
		{
			name:    "empty body should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse task data`)},
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
func getTask(t *testing.T, api http.Handler) {

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
			name:    "unknown task ID should return 404",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 1 not found`)},
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

func completeTask(t *testing.T, api http.Handler) {

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
			name:    "unknown task ID should return 404",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/task/1/complete"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 1 not found`)},
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

func clearTask(t *testing.T, api http.Handler) {

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
			name:    "unknown task ID should return 404",
			h:       api,
			args:    args{method: "DELETE", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Task ID 1 not found`)},
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

func clearCompletedTasks(t *testing.T, api http.Handler) {

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
			name:    "clearing empty task list should return 200 with 0 count",
			h:       api,
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

func listSchedules(t *testing.T, api http.Handler) {

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
			name:    "should return 200 empty list",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{}`)},
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

func addSchedule(t *testing.T, api http.Handler) {

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
			name:    "empty hourly schedule should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "hourly"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "empty/invalid schedule should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data: invalid frequency`)},
		},
		{
			name:    "invalid JSON should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{{{`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data`)},
		},
		{
			name:    "empty body should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: test.Strp(`Error: could not parse schedule data`)},
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

func getSchedule(t *testing.T, api http.Handler) {
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
			name:    "unknown ID should return 404",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 1 not found`)},
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

func pauseSchedule(t *testing.T, api http.Handler) {
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
			name:    "unknown ID should return 404",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/pause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 1 not found`)},
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
func unpauseSchedule(t *testing.T, api http.Handler) {
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
			name:    "unknown ID should return 404",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/schedule/1/unpause"},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 1 not found`)},
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

func addRecurringTask(t *testing.T, api http.Handler) {
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
			name:    "unknown ID should return 404",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusNotFound, bodyContains: test.Strp(`Schedule ID 1 not found`)},
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
