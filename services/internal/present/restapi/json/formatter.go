package json

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OutTimeFormat defines the date/time format of output fields
const OutTimeFormat = time.RFC3339Nano

// Logger interface needed for log messages
type Logger interface {
	Printf(format string, v ...interface{})
}

// ResponseFormatter defines a generic formatter interface
type ResponseFormatter interface {
	WriteResponse(w http.ResponseWriter, res []byte, statusCode int)
	WriteEmpty(w http.ResponseWriter, statusCode int)
	ErrUnauthorized(w http.ResponseWriter)
	Errorf(format string, a ...interface{}) []byte
	Error(a interface{}) []byte
}

// Formatter formats application data into JSON for output
type Formatter struct {
	l Logger
}

// NewFormatter creates a new Formatter instance
func NewFormatter(l Logger) *Formatter {
	return &Formatter{l}
}

type outError struct {
	Error string `json:"error"`
}

// WriteResponse writes the complete output response
func (f *Formatter) WriteResponse(w http.ResponseWriter, res []byte, statusCode int) {
	f.WriteEmpty(w, statusCode)
	w.Write(res)
}

// WriteEmpty writes a complete empty output response
func (f *Formatter) WriteEmpty(w http.ResponseWriter, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

// ErrUnauthorized writes an unauthorized response
func (f *Formatter) ErrUnauthorized(w http.ResponseWriter) {
	f.WriteEmpty(w, 401)
	w.Write(f.Error("Unauthorized"))
}

// Errorf formats a format string and args to JSON
func (f *Formatter) Errorf(format string, a ...interface{}) []byte {
	return f.Error(fmt.Sprintf(format, a...))
}

// Error formats an error message to JSON
func (f *Formatter) Error(a interface{}) []byte {
	outError := &outError{
		Error: fmt.Sprint(a),
	}

	o, mErr := json.Marshal(outError)
	if mErr != nil {
		f.l.Printf("problem marshalling JSON error response: %v (error struct: %v)", mErr, outError)
	}
	return o
}
