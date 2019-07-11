package postgres

import (
	"database/sql"
	"fmt"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/pqerr"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

// UserRepo handles persisting user data
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo instantiates a new UserRepo
func NewUserRepo(conn DBConn) (repo *UserRepo, err error) {
	if conn.DB == nil {
		return nil, fmt.Errorf("DB connection is nil")
	}
	return &UserRepo{db: conn.DB}, nil
}

// AddExternal adds a user and associates it to a provider and external ID
func (r *UserRepo) AddExternal(u *user.User, providerID string, externalID string) usecase.Error {

	txn, err := r.db.Begin()
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error starting DB transaction: %v", err)
	}
	defer txn.Rollback()

	id := u.ID().String()

	// Insert into user_account table
	addUserCommand := "INSERT INTO user_account (id, displayname) VALUES ($1, $2);"
	_, err = txn.Exec(addUserCommand, id, u.DisplayName())
	if err != nil {
		if pqerr.Eq(err, pqerr.UniqueViolation) {
			return usecase.NewError(usecase.ErrDuplicateRecord, "user with id %v already exists", id)
		}
		return usecase.NewError(usecase.ErrUnknown, "error inserting new user '%v': %v", id, err)
	}

	// Insert into user_external table
	addExternalCommand := "INSERT INTO user_external (user_id, provider, external_id) VALUES ($1, $2, $3);"
	_, err = txn.Exec(addExternalCommand, id, providerID, externalID)
	if err != nil {
		if pqerr.Eq(err, pqerr.UniqueViolation) {
			return usecase.NewError(usecase.ErrDuplicateRecord, "external id %v for provider %v already exists", externalID, providerID)
		}
		return usecase.NewError(usecase.ErrUnknown, "error inserting user '%v' external IDs: %v", id, err)
	}

	txn.Commit()

	return nil
}

// Update updates a user
func (r *UserRepo) Update(u *user.User) usecase.Error {

	id := u.ID().String()
	q := "UPDATE user_account SET displayname = $1 WHERE id = $2"
	res, err := r.db.Exec(q, u.DisplayName(), id)
	if err != nil {
		return usecase.NewError(usecase.ErrUnknown, "error updating user '%v': %v", id, err)
	}
	if count, _ := res.RowsAffected(); count != 1 {
		return usecase.NewError(usecase.ErrRecordNotFound, "user id '%v' not found during update", id)
	}

	return nil
}

// GetExternal gets a user given its provider and external ID
func (r *UserRepo) GetExternal(providerID string, externalID string) (*user.User, usecase.Error) {

	q := "SELECT user_account.* FROM user_account JOIN user_external ON user_account.id = user_external.user_id WHERE user_external.provider = $1 AND user_external.external_id = $2 LIMIT 1;"
	var d struct {
		id          string
		displayname string
	}
	err := r.db.QueryRow(q, providerID, externalID).Scan(&d.id, &d.displayname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.NewError(usecase.ErrRecordNotFound, "user not found by provider %v and external ID %v", providerID, externalID)
		}
		return nil, usecase.NewError(usecase.ErrUnknown, "error getting user: %v", err)
	}
	user, err := user.NewRaw(d.id, d.displayname)
	if err != nil {
		return nil, usecase.NewError(usecase.ErrUnknown, "error parsing user data: %v", err)
	}

	return user, nil
}
