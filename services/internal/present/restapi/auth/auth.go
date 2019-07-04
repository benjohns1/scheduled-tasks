package auth

import (
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

/*
	@TODO: implement Auth0
	"gopkg.in/square/go-jose.v2"
	"github.com/auth0-community/go-auth0"
*/

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

var signingKey = []byte("asdfasdfasdf")

// GetToken returns a JWT token
func GetToken(l Logger) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		claims := &jwt.StandardClaims{
			Subject:   "lando.calrissian@email.com.invalid",
			Issuer:    "Scheduled Tasks Test App",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

		tokenString, err := token.SignedString(signingKey)
		if err != nil {
			l.Printf("error signing token: %v", err)
			return
		}

		w.Write([]byte(tokenString))
	}
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	},
	SigningMethod: jwt.SigningMethodHS512,
})

// Handler httprouter middleware that handles JWT request auth
func Handler(next httprouter.Handle) httprouter.Handle {
	return next
	//return httprouterwrap.Wrap(jwtMiddleware.Handler(httprouterwrap.Unwrap(next)))
}
