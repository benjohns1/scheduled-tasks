package restapi

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
	mapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	scheduleapi "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/schedule"
	taskapi "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/task"
	userapi "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// New creates a REST API server
func New(l Logger, a auth.Authenticator, checkSchedule chan<- bool, userRepo usecase.UserRepo, taskRepo usecase.TaskRepo, scheduleRepo usecase.ScheduleRepo) (api http.Handler) {

	r := httprouter.New()
	f := mapper.NewFormatter(l)
	a.SetFormatter(f)
	prefix := "/api/v1"
	taskapi.Handle(r, prefix, l, f, taskRepo)
	scheduleapi.Handle(r, prefix, l, f, checkSchedule, scheduleRepo)
	userapi.Handle(r, prefix, l, f, userRepo)

	r.HandleMethodNotAllowed = false
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.WriteResponse(w, f.Error("Not found"), 404)
	})
	return a.Authenticate(auth.HydrateUser(userRepo, l, f, false, r))
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
