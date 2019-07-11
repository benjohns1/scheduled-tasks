package json

import (
	"encoding/json"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	format "github.com/benjohns1/scheduled-tasks/services/internal/present/restapi/json"
)

// Formatter formats application data into JSON for output
type Formatter struct {
	format.ResponseFormatter
}

// NewFormatter creates a new Formatter instance
func NewFormatter(rf format.ResponseFormatter) *Formatter {
	return &Formatter{rf}
}

type outUserID struct {
	ID user.ID `json:"id"`
}

// UserID formats a UserID to JSON
func (f *Formatter) UserID(id user.ID) ([]byte, error) {
	o := &outUserID{
		ID: id,
	}
	return json.Marshal(o)
}
