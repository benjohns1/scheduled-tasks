package httprouterwrap

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type httprouterwrap struct {
	http.ResponseWriter
	Params httprouter.Params
}

// Wrap wraps a standard net/http middleware handler and returns a julienschmidt/httprouter middleware handler
func Wrap(next http.Handler) httprouter.Handle {
	return WrapFunc(next.ServeHTTP)
}

// Unwrap unwraps a wrapped julienschmidt/httprouter middleware handler and returns a standard net/http middleware handler
func Unwrap(next httprouter.Handle) http.Handler {
	return http.HandlerFunc(UnwrapFunc(next))
}

// WrapFunc wraps a standard net/http middleware handler function and returns a julienschmidt/httprouter middleware handler
func WrapFunc(next http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		wrapped := &httprouterwrap{ResponseWriter: w, Params: ps}
		next(wrapped, r)
	}
}

// UnwrapFunc unwraps a wrapped julienschmidt/httprouter middleware handler and returns a standard net/http middleware handler function
func UnwrapFunc(next httprouter.Handle) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if wrapped, ok := w.(*httprouterwrap); ok {
			next(wrapped.ResponseWriter, r, wrapped.Params)
		} else {
			next(w, r, httprouter.Params{})
		}
	}
}
