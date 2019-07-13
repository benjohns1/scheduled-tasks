package auth

import (
	"net/http"
	"strings"

	"github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// Auth0 contains dependencies for the Auth0 handler
type Auth0 struct {
	Auth
	c Auth0Config
}

// Auth0Config contains configuration options for Auth0 handler
type Auth0Config struct {
	Secret   []byte
	Audience []string
	Domain   string
}

// NewAuth0 returns a new Auth struct
func NewAuth0(l Logger, c Auth0Config) *Auth0 {
	return &Auth0{Auth{l: l}, c}
}

type claims struct {
	Issuer  string `json:"iss"`
	Subject string `json:"sub"`
	Scope   string `json:"scope"`
}

// Authenticate authenticates a request and calls the next handler
func (a *Auth0) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretProvider := auth0.NewKeyProvider(a.c.Secret)
		configuration := auth0.NewConfiguration(secretProvider, a.c.Audience, "https://"+a.c.Domain+"/", jose.HS256)
		validator := auth0.NewValidator(configuration, nil)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			a.l.Printf("Error parsing token:", err)
			a.l.Printf("Token is not valid:", token)
			if a.f != nil {
				a.f.WriteResponse(w, a.f.Error("Unauthorized"), http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}
			return
		}
		c := claims{}
		validator.Claims(r, token, &c)
		ps := getDefaultPerms(&c)

		next.ServeHTTP(ResponseContext{w, Context{c.Issuer, c.Subject, ps}}, r)
	})
}

func getDefaultPerms(c *claims) []Permission {
	scopes := strings.Split(c.Scope, " ")
	for _, scope := range scopes {
		switch scope {
		case "type:anon":
			return []Permission{PermNone}
		}
	}
	return GetDefaultUserPerms()
}
