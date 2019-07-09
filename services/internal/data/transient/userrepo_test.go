package transient

import (
	"reflect"
	"testing"

	"github.com/benjohns1/scheduled-tasks/services/internal/core/user"
	"github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func newUserRaw(t *testing.T, idStr string, displayname string) *user.User {
	u, err := user.NewRaw(idStr, displayname)
	if err != nil {
		t.Fatalf("Error creating new raw user: %v", err)
	}
	return u
}

func TestUserRepo_Add(t *testing.T) {
	r := NewUserRepo()
	emptyUser := user.New("")
	basicUser := user.New("Test Displayname")
	dupeUser1 := newUserRaw(t, "111e1111-e89b-12d3-a456-426655440000", "displayname1")
	dupeUser2 := newUserRaw(t, "111e1111-e89b-12d3-a456-426655440000", "displayname2")

	if err := r.Add(dupeUser1, "p1", "e1"); err != nil {
		t.Fatalf("Error adding user")
	}

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
			name:    "should add an empty user",
			r:       r,
			args:    args{emptyUser, "", ""},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should add a basic user",
			r:       r,
			args:    args{basicUser, "p2", "e2"},
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should return duplicate error",
			r:       r,
			args:    args{dupeUser2, "p1", "e1"},
			wantErr: usecase.ErrDuplicateRecord,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Add(tt.args.u, tt.args.providerID, tt.args.externalID)
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UserRepo.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepo_Update(t *testing.T) {
	r := NewUserRepo()
	u := user.New("display name")
	r.Add(u, "", "")
	u.UpdateDisplayName("new display name")

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
			name:    "should not return error code",
			r:       r,
			args:    args{u},
			wantErr: usecase.ErrNone,
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
	r := NewUserRepo()
	u := user.New("display name")
	r.Add(u, "p1", "e1")

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
			name:    "should successfully get a user",
			r:       r,
			args:    args{"p1", "e1"},
			want:    u,
			wantErr: usecase.ErrNone,
		},
		{
			name:    "should return an error",
			r:       r,
			args:    args{"p1", "invalid-external-id"},
			want:    nil,
			wantErr: usecase.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetExternal(tt.args.providerID, tt.args.externalID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserRepo.GetExternal() got = %v, want %v", got, tt.want)
			}
			if ((err == nil) != (tt.wantErr == usecase.ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("UserRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
