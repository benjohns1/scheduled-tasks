// +build integration

package restapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benjohns1/scheduled-tasks/internal/data/transient"
	"github.com/julienschmidt/httprouter"

	mapper "github.com/benjohns1/scheduled-tasks/internal/present/restapi/json"
)

type LoggerMock struct{}

func (l *LoggerMock) Printf(format string, v ...interface{}) {}

func Test_listSchedules(t *testing.T) {

	scheduleRepo := transient.NewScheduleRepo()
	l := &LoggerMock{}
	f := mapper.NewFormatter(l)
	r := httprouter.New()
	r.GET("/schedule", listSchedules(l, f, scheduleRepo))

	type args struct {
		method string
		url    string
		body   io.Reader
	}
	type asserts struct {
		statusEquals int
		bodyEquals   string
	}
	tests := []struct {
		name    string
		h       http.Handler
		args    args
		asserts asserts
	}{
		{
			name:    "should return 200 empty list",
			h:       r,
			args:    args{method: "GET", url: "/schedule"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: `{}`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req)
			if rr.Code != tt.asserts.statusEquals {
				t.Errorf("listSchedules() status code = %v, want %v", rr.Code, tt.asserts.statusEquals)
			}
			if rr.Body.String() != tt.asserts.bodyEquals {
				t.Errorf("listSchedules() response body = %v, want %v", rr.Body.String(), tt.asserts.bodyEquals)
			}
		})
	}
}
