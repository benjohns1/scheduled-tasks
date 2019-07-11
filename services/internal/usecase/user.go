package usecase

import (
	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
)

// UserRepo defines the user repository interface required by use cases
type UserRepo interface {
	AddExternal(u *user.User, providerID string, externalID string) Error
	Update(*user.User) Error
	GetExternal(providerID string, externalID string) (*user.User, Error)
}

// AddOrUpdateExternalUser looks up a user by an external provider's ID, then either adds them or updates their displayname if needed
func AddOrUpdateExternalUser(r UserRepo, providerID string, externalID string, displayname string) (*user.User, Error) {
	u, err := r.GetExternal(providerID, externalID)
	if err != nil {
		if err.Code() != ErrRecordNotFound {
			return nil, err
		}

		u = user.New(displayname)
		if err := r.AddExternal(u, providerID, externalID); err != nil {
			return nil, err.Prefix("error adding external user: ")
		}
		return u, nil
	}

	if u.DisplayName() != displayname {
		u.UpdateDisplayName(displayname)
		if err := r.Update(u); err != nil {
			return nil, err.Prefix("error updating user: ")
		}
	}
	return u, nil
}
