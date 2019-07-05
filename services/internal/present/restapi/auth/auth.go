package auth

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// Formatter defines the formatter interface for output responses
type Formatter interface {
	WriteResponse(w http.ResponseWriter, res []byte, statusCode int)
	Error(a interface{}) []byte
}

// Auth base struct for auth implementations
type Auth struct {
	l Logger
	f Formatter
	Authorizer
}

func New(l Logger) *Auth {
	return &Auth{l: l}
}

// SetFormatter sets the formatter on the Auth struct
func (a *Auth) SetFormatter(f Formatter) {
	a.f = f
}

// Handle dummy method for auth handler
func (a *Auth) Handle(next httprouter.Handle) httprouter.Handle {
	return next
}

// Authorizer interface for authorization
type Authorizer interface {
	SetFormatter(f Formatter)
	Handle(next httprouter.Handle) httprouter.Handle
}
