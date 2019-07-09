package user

// User a user of the scheduled task system
type User struct {
	id          ID
	displayname string
}

// New creates a new user entity
func New(displayname string) *User {
	id := NewID()
	return &User{id, displayname}
}

// NewRaw instantiates a user entity with all available fields
func NewRaw(idStr string, displayname string) (*User, error) {
	id, err := ParseID(idStr)
	if err != nil {
		return nil, err
	}
	return &User{id, displayname}, nil
}

// ID returns the user's unique ID
func (u *User) ID() ID {
	return u.id
}

// DisplayName returns the user's display name
func (u *User) DisplayName() string {
	return u.displayname
}

// UpdateDisplayName updates the user's display name
func (u *User) UpdateDisplayName(displayname string) {
	u.displayname = displayname
}
