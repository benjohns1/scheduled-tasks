package restapi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	scheduleapi "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/schedule"
	taskapi "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/task"
	mapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// New creates a REST API server
func New(l Logger, checkSchedule chan<- bool, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo) (api http.Handler) {

	r := httprouter.New()
	f := mapper.NewFormatter(l)
	prefix := "/api/v1"
	taskapi.Handle(r, prefix, l, f, taskRepo)
	scheduleapi.Handle(r, prefix, l, f, checkSchedule, scheduleRepo)

	return r
}

// Serve starts an API server
func Serve(l Logger, api http.Handler) (closed <-chan bool) {

	onClosed := make(chan bool)

	// Start API server
	go func() {
		port := 8080
		if val, err := strconv.Atoi(os.Getenv("APPLICATION_PORT")); err == nil {
			port = val
		}

		l.Printf("starting server on port %d", port)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), api)
		if err != nil {
			l.Printf("http server error: %v", err)
		}
		l.Printf("server exiting")
		onClosed <- true
	}()

	return onClosed
}