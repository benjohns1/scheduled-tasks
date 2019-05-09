// +build integration

package restapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/benjohns1/scheduled-tasks/internal/data/transient"
)

type LoggerMock struct{}

func (l *LoggerMock) Printf(format string, v ...interface{}) {}

func strp(str string) *string {
	return &str
}

func Test_listSchedules(t *testing.T) {

	taskRepo := transient.NewTaskRepo()
	scheduleRepo := transient.NewScheduleRepo()
	l := &LoggerMock{}
	
	api := New(l, taskRepo, scheduleRepo)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains   *string
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
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: strp(`{}`)},
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

func Test_addSchedule(t *testing.T) {

	taskRepo := transient.NewTaskRepo()
	scheduleRepo := transient.NewScheduleRepo()
	l := &LoggerMock{}
	
	api := New(l, taskRepo, scheduleRepo)

	type args struct {
		method string
		url    string
		body   string
	}
	type asserts struct {
		statusEquals int
		bodyEquals   *string
		bodyContains   *string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "empty schedule should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: strp(`{"id":1}`)},
		},
		{
			name:    "invalid JSON should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: `{{{`},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: strp(`Error: could not parse schedule data`)},
		},
		{
			name:    "empty body should return 400",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/schedule/", body: ``},
			asserts: asserts{statusEquals: http.StatusBadRequest, bodyContains: strp(`Error: could not parse schedule data`)},
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