package test

import (
	"net/http"

	"github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/auth"
)

// AuthMock contains dependencies for the AuthMock handler
type AuthMock struct {
	auth.Auth
}

// MockClaims contains static claims to inject in request
type MockClaims struct {
	Issuer      string
	Subject     string
	Permissions []auth.Permission
}

type injectContext struct {
	http.ResponseWriter
	claims MockClaims
}

// NewAuthMock returns a mock AuthMock struct
func NewAuthMock(l auth.Logger) *AuthMock {
	return &AuthMock{*auth.New(l)}
}

// Authenticate authenticates a request and calls the next handler
func (a *AuthMock) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ic, ok := w.(injectContext); ok {
			next.ServeHTTP(auth.ResponseContext{ResponseWriter: w, Auth: auth.Context(ic.claims)}, r)
			return
		}
		next.ServeHTTP(auth.ResponseContext{ResponseWriter: w, Auth: auth.Context{}}, r)
	})
}

// InjectClaims injects claims into the current request
func InjectClaims(claims MockClaims, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(injectContext{w, claims}, r)
	})
}
