package user

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	responseMapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
	mapper "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/user/json"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter defines the formatter interface for output responses
type Formatter interface {
	responseMapper.ResponseFormatter
}

// Parser defines the parser interface for parsing input requests
type Parser interface {
	AddOrUpdateUser(b io.Reader) (mapper.User, error)
}

// Auth defines the auth interface for auth middleware
type Auth interface {
	Handle(next httprouter.Handle) httprouter.Handle
}

// Handle adds user handling endpoints
func Handle(r *httprouter.Router, a Auth, prefix string, l Logger, rf responseMapper.ResponseFormatter, userRepo usecase.UserRepo) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(rf)

	pre := prefix + "/user"
	r.PUT(pre+"/external/:providerID/:userID/addOrUpdate", a.Handle(addOrUpdateExternalUser(l, p, f, userRepo)))
}

func addOrUpdateExternalUser(l Logger, p Parser, f Formatter, userRepo usecase.UserRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		providerID := ps.ByName("providerID")
		userID := ps.ByName("userID")
		if providerID == "" || userID == "" {
			l.Printf("valid provider and user IDs required")
			f.WriteResponse(w, f.Error("Error: valid provider and user IDs required"), 404)
			return
		}
		userData, ucerr := p.AddOrUpdateUser(r.Body)
		defer r.Body.Close()
		if ucerr != nil {
			l.Printf("error parsing AddOrUpdateUser data: %v", ucerr)
			f.WriteResponse(w, f.Errorf("Error: could not parse user data: %v", ucerr), 400)
			return
		}

		_, ucerr = usecase.AddOrUpdateExternalUser(userRepo, providerID, userID, userData.DisplayName)
		if ucerr != nil {
			l.Printf("error adding or updating external user: %v", ucerr)
			f.WriteResponse(w, f.Error("Error adding or updating external user"), 500)
			return
		}
		f.WriteEmpty(w, 204)
	}
}
