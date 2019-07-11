package transient

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// UserRepo maintains an in-memory cache of users
type UserRepo struct {
	users    map[user.ID]*user.User
	external map[providerKey]user.ID
}

type providerKey struct {
	provider string
	id       string
}

// NewUserRepo instantiates a new UserRepo
func NewUserRepo() *UserRepo {
	return &UserRepo{users: make(map[user.ID]*user.User), external: make(map[providerKey]user.ID)}
}

// AddExternal adds a user to the memory cache
func (r *UserRepo) AddExternal(u *user.User, providerID string, externalID string) usecase.Error {
	id := u.ID()
	if (id == user.ID{}) {
		return usecase.NewError(usecase.ErrInvalidID, "user ID cannot be empty when adding to repo")
	}
	if _, exists := r.users[id]; exists {
		return usecase.NewError(usecase.ErrDuplicateRecord, "user ID %v already exists in repo", id)
	}
	provider := providerKey{providerID, externalID}
	if mappedUserID, exists := r.external[provider]; exists {
		return usecase.NewError(usecase.ErrDuplicateRecord, "external provider ID %v already exists in repo for user ID %v", provider, mappedUserID)
	}
	r.users[id] = u
	r.external[provider] = id
	return nil
}

// Update updates a user
func (r *UserRepo) Update(u *user.User) usecase.Error {
	id := u.ID()
	if _, ok := r.users[id]; !ok {
		return usecase.NewError(usecase.ErrRecordNotFound, "no user with ID %v", id)
	}

	r.users[id] = u

	return nil
}

// GetExternal gets a user given its provider and external ID
func (r *UserRepo) GetExternal(providerID string, externalID string) (*user.User, usecase.Error) {

	provider := providerKey{providerID, externalID}
	id, ok := r.external[provider]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no user with provider ID %v", provider)
	}

	u, ok := r.users[id]
	if !ok {
		return nil, usecase.NewError(usecase.ErrRecordNotFound, "no user with mapped ID: %v found given provider ID: %v", id, provider)
	}
	return u, nil
}
