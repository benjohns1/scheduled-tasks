package auth

import (
	"fmt"
	"net/http"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
	"github.com/julienschmidt/httprouter"
)

// UserContext wraps http.ResponseWriter in a context that provides a hydrated domain user to other handlers
type UserContext struct {
	http.ResponseWriter
	User *user.User
	Auth Context
}

// HydrateUser middleware hydrates a UserContext with a user
// will respond with a 401 unauthorized response if required is set to true and no user could be found
func HydrateUser(userRepo usecase.UserRepo, l Logger, f Formatter, required bool, autoCreate bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, a, ok := hydrateUser(w, userRepo, l, f, required, autoCreate); ok {
			next.ServeHTTP(UserContext{w, u, a}, r)
		}
	})
}

// HRHydrateUser wraps HydrateUser in httprouter middleware
func HRHydrateUser(userRepo usecase.UserRepo, l Logger, f Formatter, required bool, autoCreate bool, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if u, a, ok := hydrateUser(w, userRepo, l, f, required, autoCreate); ok {
			next(UserContext{w, u, a}, r, ps)
		}
	}
}

func hydrateUser(w http.ResponseWriter, userRepo usecase.UserRepo, l Logger, f Formatter, required bool, autoCreate bool) (*user.User, Context, bool) {
	a, ok := w.(ResponseContext)
	if !ok {
		l.Printf("Invalid auth context, could not parse user: %v", w)
		f.WriteResponse(w, f.Error("Error parsing user"), 500)
		return nil, Context{}, false
	}
	c := Context(a.Auth)

	u, err := usecase.GetExternalUser(userRepo, a.Auth.Issuer, a.Auth.Subject)

	if err != nil {
		if err.Code() != usecase.ErrRecordNotFound {
			l.Printf("Error finding user: %v", err)
			f.WriteResponse(w, f.Error("Error finding user"), 500)
			return nil, c, false
		}
		if required {
			l.Printf("Error finding authorized user from token: %v", err)
			f.WriteResponse(w, f.Error("Could not find authorized user from token"), 401)
			return nil, c, false
		}
		u = &user.User{}
	}

	return u, c, true
}

// FormatProvider formats a provider string from a request to the DB and issuer format
func FormatProvider(provider string) string {
	return fmt.Sprintf("https://%v/", provider)
}
