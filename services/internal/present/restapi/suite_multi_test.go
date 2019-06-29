// +build integration

package restapi_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/clock"
	format "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/test"
)

func TestTransientRESTAPIMulti(t *testing.T) {
	tester := test.NewTransientTester()
	defer tester.Close()
	suiteMulti(t, tester)
}

func TestPostgresRESTAPIMulti(t *testing.T) {
	tester := test.NewPostgresTester()
	defer tester.Close()
	suiteMulti(t, tester)
}

func suiteMulti(t *testing.T, tester test.Tester) {
	addListGetCompleteTasks(t, tester.NewAPI())
	addListGetSchedules(t, tester.NewAPI())
	addRecurringTasksToEmptySchedule(t, tester.NewAPI())
	addRemoveListSchedule(t, tester.NewAPI())
}

func addListGetCompleteTasks(t *testing.T, api http.Handler) {

	now := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	nowStr := now.Format(format.OutTimeFormat)
	prevClock := clock.Get()
	clockMock := clock.NewStaticMock(now)
	clock.Set(clockMock)
	defer clock.Set(prevClock)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals    int
		bodyEquals      *string
		bodyContains    *string
		bodyNotEquals   *string
		bodyNotContains *string
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
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(fmt.Sprintf(`{"1":{"id":1,"name":"task1","description":"task1 description","completedTime":null,"createdTime":"%v"},"2":{"id":2,"name":"task2","description":"task2 description","completedTime":null,"createdTime":"%v"},"3":{"id":3,"name":"task3","description":"task3 description","completedTime":null,"createdTime":"%v"}}`, nowStr, nowStr, nowStr))},
		},
		{
			name:    "get task ID 1 should return incompleted task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(fmt.Sprintf(`{"id":1,"name":"task1","description":"task1 description","completedTime":null,"createdTime":"%v"}`, nowStr))},
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
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(fmt.Sprintf(`{"id":1,"name":"task1","description":"task1 description","completedTime":"%v","createdTime":"%v"}`, nowStr, nowStr))},
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

func addRemoveListSchedule(t *testing.T, api http.Handler) {
	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals    int
		bodyEquals      *string
		bodyContains    *string
		bodyNotEquals   *string
		bodyNotContains *string
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
			name:    "new hourly schedule should return 201 and ID 1",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Hour", "atMinutes":[0]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "new hourly schedule with tasks should return 201 and ID 2",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0,15,30], "tasks": [{"name":"rtask1","description":"rtask1 desc"}]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "new hourly schedule with interval and offset should return 201 and ID 3",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0], "interval": 2, "offset": 1}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":3}`)},
		},
		{
			name:    "get schedule ID 1 should return hourly schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 2 should return hourly schedule with 1 recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":2,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,15,30],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]}`)},
		},
		{
			name:    "get schedule ID 3 should return hourly schedule with no recurring tasks and with interval and offset",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/3"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":3,"frequency":"Hour","interval":2,"offset":1,"atMinutes":[0],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "removing schedule 1 should return 204",
			h:       api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "removing schedule 3 should return 204",
			h:       api,
			args:    args{method: "DELETE", url: "/api/v1/schedule/3"},
			asserts: asserts{statusEquals: http.StatusNoContent},
		},
		{
			name:    "list return 200 list with 1 schedule with ID 2",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"2":{"id":2,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,15,30],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]}}`)},
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

func addListGetSchedules(t *testing.T, api http.Handler) {
	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals    int
		bodyEquals      *string
		bodyContains    *string
		bodyNotEquals   *string
		bodyNotContains *string
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
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Hour", "atMinutes":[]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "new hourly schedule should return 201 and ID 2",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0,30], "paused":true}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":2}`)},
		},
		{
			name:    "new hourly schedule should return 201 and ID 3",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0,30,59], "tasks": [{"name": "rtask1", "description": "rtask1 desc"}]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":3}`)},
		},
		{
			name:    "new hourly schedule with interval and offset should return 201 and ID 4",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0], "interval": 2, "offset": 1}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":4}`)},
		},
		{
			name:    "day schedule should return 201 and ID 5",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Day", "atMinutes":[0,30], "atHours":[3,6]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":5}`)},
		},
		{
			name:    "week schedule should return 201 and ID 6",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Week", "atMinutes":[0,30], "atHours":[3,6], "onDaysOfWeek":["Wednesday","Thursday"]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":6}`)},
		},
		{
			name:    "month schedule should return 201 and ID 7",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Month", "atMinutes":[15], "atHours":[1], "onDaysOfMonth":[1,15,31]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":7}`)},
		},
		{
			name:    "new hourly schedule with invalid ranges for interval and offset should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency": "Hour", "atMinutes": [0], "interval": 0, "offset": -1}`},
			asserts: asserts{statusEquals: http.StatusBadRequest},
		},
		{
			name:    "get schedule ID 1 should return empty schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 2 should return hourly schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/2"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":2,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,30],"paused":true,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 3 should return hourly schedule with 1 recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/3"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":3,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,30,59],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]}`)},
		},
		{
			name:    "get schedule ID 4 should return empty schedule with no recurring tasks and interval and offset",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/4"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":4,"frequency":"Hour","interval":2,"offset":1,"atMinutes":[0],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 5 should return day schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/5"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":5,"frequency":"Day","interval":1,"offset":0,"atMinutes":[0,30],"atHours":[3,6],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 6 should return week schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/6"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":6,"frequency":"Week","interval":1,"offset":0,"atMinutes":[0,30],"atHours":[3,6],"onDaysOfWeek":["Wednesday","Thursday"],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "get schedule ID 7 should return month schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/7"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":7,"frequency":"Month","interval":1,"offset":0,"atMinutes":[15],"atHours":[1],"onDaysOfMonth":[1,15,31],"paused":false,"tasks":[]}`)},
		},
		{
			name:    "should return 200 list with 7 schedules",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"1":{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[]},"2":{"id":2,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,30],"paused":true,"tasks":[]},"3":{"id":3,"frequency":"Hour","interval":1,"offset":0,"atMinutes":[0,30,59],"paused":false,"tasks":[{"name":"rtask1","description":"rtask1 desc"}]},"4":{"id":4,"frequency":"Hour","interval":2,"offset":1,"atMinutes":[0],"paused":false,"tasks":[]},"5":{"id":5,"frequency":"Day","interval":1,"offset":0,"atMinutes":[0,30],"atHours":[3,6],"paused":false,"tasks":[]},"6":{"id":6,"frequency":"Week","interval":1,"offset":0,"atMinutes":[0,30],"atHours":[3,6],"onDaysOfWeek":["Wednesday","Thursday"],"paused":false,"tasks":[]},"7":{"id":7,"frequency":"Month","interval":1,"offset":0,"atMinutes":[15],"atHours":[1],"onDaysOfMonth":[1,15,31],"paused":false,"tasks":[]}}`)},
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

func addRecurringTasksToEmptySchedule(t *testing.T, api http.Handler) {
	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals    int
		bodyEquals      *string
		bodyContains    *string
		bodyNotEquals   *string
		bodyNotContains *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "empty hourly schedule should return 201 and ID 1",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{"frequency":"Hour", "atMinutes":[]}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(`{"id":1}`)},
		},
		{
			name:    "get schedule ID 1 should return empty schedule with no recurring tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[]}`)},
		},
		{
			name:    "should return 200 list with 1 schedule",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"1":{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[]}}`)},
		},
		{
			name:    "adding recurring task to schedule ID 1 should return 201",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/1/task/", body: `{"name":"task1","description":"task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: test.Strp(``)},
		},
		{
			name:    "get schedule ID 1 should return schedule with 1 task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[{"name":"task1","description":"task1 description"}]}`)},
		},
		{
			name:    "should return 200 list with 1 schedule with 1 task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/schedule/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: test.Strp(`{"1":{"id":1,"frequency":"Hour","interval":1,"offset":0,"paused":false,"tasks":[{"name":"task1","description":"task1 description"}]}}`)},
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
