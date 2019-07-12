package user

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
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

// Handle adds user handling endpoints
func Handle(r *httprouter.Router, prefix string, l Logger, rf responseMapper.ResponseFormatter, userRepo usecase.UserRepo) {

	p := mapper.NewParser()
	f := mapper.NewFormatter(rf)

	pre := prefix + "/user"
	r.PUT(pre+"/external/:providerID/:userID/addOrUpdate", authorize(l, f, addOrUpdateExternalUser(l, p, f, userRepo)))
}

func authorize(l Logger, f Formatter, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userContext, ok := w.(auth.UserContext)
		if !ok {
			l.Printf("invalid authorization context from http.ResponseWriter: %v", w)
			f.WriteResponse(w, f.Error("Internal authorization error"), 500)
			return
		}

		// Enforce users may only create/update themselves
		if userContext.Auth.HasPerm(auth.PermUpsertUserSelf) {
			providerID := ps.ByName("providerID")
			userID := ps.ByName("userID")
			if auth.FormatProvider(providerID) != userContext.Auth.Issuer || userID != userContext.Auth.Subject {
				l.Printf("external user credentials (%v, %v) do not match authorization token: %v", providerID, userID, userContext)
				f.WriteResponse(w, f.Error("Not authorized"), 401)
				return
			}
			next(w, r, ps)
			return
		}

		f.ErrUnauthorized(w)
	}
}

func addOrUpdateExternalUser(l Logger, p Parser, f Formatter, userRepo usecase.UserRepo) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Validate URL params
		providerID := ps.ByName("providerID")
		userID := ps.ByName("userID")
		if providerID == "" || userID == "" {
			l.Printf("valid provider and user IDs required")
			f.WriteResponse(w, f.Error("Error: valid provider and user IDs required"), 404)
			return
		}
		provider := auth.FormatProvider(providerID)

		// Validate URL params match logged-in user
		u1 := auth.GetUser(w)
		u2, err := usecase.GetExternalUser(userRepo, provider, userID)
		if err != nil || !u1.ID().Equals(u2.ID()) {
			l.Printf("logged-in user doesn't match URL specified user (%v != %v), err: %v", u1, u2, err)
			f.WriteResponse(w, f.Errorf("Error: logged-in user doesn't match URL specified user"), 400)
			return
		}

		userData, ucerr := p.AddOrUpdateUser(r.Body)
		defer r.Body.Close()
		if ucerr != nil {
			l.Printf("error parsing AddOrUpdateUser data: %v", ucerr)
			f.WriteResponse(w, f.Errorf("Error: could not parse user data: %v", ucerr), 400)
			return
		}

		_, ucerr = usecase.AddOrUpdateExternalUser(userRepo, provider, userID, userData.DisplayName)
		if ucerr != nil {
			l.Printf("error adding or updating external user: %v", ucerr)
			f.WriteResponse(w, f.Error("Error adding or updating external user"), 500)
			return
		}
		f.WriteEmpty(w, 204)
	}
}
