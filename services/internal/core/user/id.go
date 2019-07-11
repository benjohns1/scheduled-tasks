package user

import (
	"github.com/google/uuid"
)

// ID unique user identifier
type ID struct {
	id uuid.UUID
}

// NewID generates a new ID
func NewID() ID {
	return ID{id: uuid.New()}
}

// ParseID creates an ID from a preexisting string value
func ParseID(val string) (ID, error) {
	id, err := uuid.Parse(val)
	if err != nil {
		return ID{}, err
	}
	return ID{id: id}, nil
}

// Equals determines if two user IDs are equal
func (val ID) Equals(other interface{}) bool {
	if otherVal, ok := other.(ID); ok {
		return val.id == otherVal.id
	}
	return false
}

// String returns the string representation of the ID
func (val ID) String() string {
	return val.id.String()
}

// StringPtr returns a pointer to the string representation of the ID, or nil if it is the zero-value
func (val ID) StringPtr() *string {
	if (val.id == uuid.UUID{}) {
		return nil
	}
	str := val.id.String()
	return &str
}
