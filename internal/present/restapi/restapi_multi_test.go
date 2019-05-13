// +build integration

package restapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_addListGetCompleteTasks(t *testing.T) {

	api := mockAPI()

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
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: strp(`{}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task1", "description": "task1 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: strp(`{"id":1}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task2", "description": "task2 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: strp(`{"id":2}`)},
		},
		{
			name:    "task with name and description should return 201 and ID",
			h:       api,
			args:    args{method: "POST", url: "/api/v1/task/", body: `{"name": "task3", "description": "task3 description"}`},
			asserts: asserts{statusEquals: http.StatusCreated, bodyEquals: strp(`{"id":3}`)},
		},
		{
			name:    "should return 200 list with 3 tasks",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: strp(`{"1":{"id":1,"name":"task1","description":"task1 description","completedTime":null},"2":{"id":2,"name":"task2","description":"task2 description","completedTime":null},"3":{"id":3,"name":"task3","description":"task3 description","completedTime":null}}`)},
		},
		{
			name:    "get task ID 1 should return incompleted task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyEquals: strp(`{"id":1,"name":"task1","description":"task1 description","completedTime":null}`)},
		},
		{
			name:    "complete ID 1 should return 204",
			h:       api,
			args:    args{method: "PUT", url: "/api/v1/task/1/complete"},
			asserts: asserts{statusEquals: http.StatusNoContent, bodyContains: strp(``)},
		},
		{
			name:    "get task ID 1 should return completed task",
			h:       api,
			args:    args{method: "GET", url: "/api/v1/task/1"},
			asserts: asserts{statusEquals: http.StatusOK, bodyContains: strp(`{"id":1,"name":"task1","description":"task1 description","completedTime":`), bodyNotContains: strp(`"completedTime":null`)},
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