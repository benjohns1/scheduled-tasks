// +build integration

package postgres_test

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres"
	. "github.com/benjohns1/scheduled-tasks/services/internal/data/postgres/test"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestNewUserRepo(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	type args struct {
		conn DBConn
	}
	tests := []struct {
		name      string
		args      args
		wantValue bool
		wantErr   bool
	}{
		{
			name:      "should not return an error",
			args:      args{conn},
			wantValue: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUserRepo(tt.args.conn)
			if (got != nil) != tt.wantValue {
				t.Errorf("NewUserRepo() = %v, wantValue %v", got, tt.wantValue)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_AddExternal(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	r, err := NewUserRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	u1 := user.New("new user display name")
	u2 := user.New("new user display name")

	type args struct {
		u          *user.User
		providerID string
		externalID string
	}
	tests := []struct {
		name    string
		r       *UserRepo
		args    args
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should add a user without error",
			r:       r,
			args:    args{u1, "provider1", "extid1"},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "duplicate provider/id should return dulicate error",
			r:       r,
			args:    args{user.New("another user display name"), "provider1", "extid1"},
			wantErr: usecase.ErrDuplicateRecord,
		},
		{
			name:    "duplicate user should return duplicate error",
			r:       r,
			args:    args{u1, "provider2", "extid2"},
			wantErr: usecase.ErrDuplicateRecord,
		},
		{
			name:    "different user with external key but the same data should be added without error",
			r:       r,
			args:    args{u2, "provider1", "extid3"},
			wantErr: usecase.ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.AddExternal(tt.args.u, tt.args.providerID, tt.args.externalID)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UserRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	r, err := NewUserRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	u1 := user.New("new user display name")
	r.AddExternal(u1, "p1", "e1")
	u1.UpdateDisplayName("updated display name")

	type args struct {
		u *user.User
	}
	tests := []struct {
		name    string
		r       *UserRepo
		args    args
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should update a user without error",
			r:       r,
			args:    args{u1},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "updating a non-existent user should return not found error",
			r:       r,
			args:    args{user.New("non-existent user")},
			wantErr: usecase.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Update(tt.args.u)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UserRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_GetExternal(t *testing.T) {
	conn, err := NewTestDBConn(DBTest)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	r, err := NewUserRepo(conn)
	if err != nil {
		t.Fatal(err)
	}

	u1 := user.New("new user display name")
	r.AddExternal(u1, "p1", "e1")

	type args struct {
		providerID string
		externalID string
	}
	tests := []struct {
		name    string
		r       *UserRepo
		args    args
		want    *user.User
		wantErr usecase.ErrorCode
	}{
		{
			name:    "should get valid user",
			r:       r,
			args:    args{"p1", "e1"},
			want:    u1,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should return not found error",
			r:       r,
			args:    args{"p1", "non-existent user"},
			want:    nil,
			wantErr: usecase.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetExternal(tt.args.providerID, tt.args.externalID)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UserRepo.GetExternal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.GetExternal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
