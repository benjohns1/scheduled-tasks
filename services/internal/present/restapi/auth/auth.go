package auth

import (
	"net/http"
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
	Authenticator
}

// ResponseContext wraps http.ResponseWriter in a context that provides authorization details to other handlers
type ResponseContext struct {
	http.ResponseWriter
	Auth Context
}

// Context contains relevant auth data from request
type Context struct {
	Issuer      string
	Subject     string
	Permissions []string
	Scope       string
}

// New creates a new base Auth struct (useful for test stubbing)
func New(l Logger) *Auth {
	return &Auth{l: l}
}

// SetFormatter sets the formatter on the Auth struct
func (a *Auth) SetFormatter(f Formatter) {
	a.f = f
}

// Authenticate stub authentication method
func (a *Auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// Authenticator interface for authorization
type Authenticator interface {
	SetFormatter(f Formatter)
	Authenticate(next http.Handler) http.Handler
}
