package user

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		displayname string
	}
	tests := []struct {
		name        string
		args        args
		wantValidID bool
		want        *User
	}{
		{
			name:        "should create a valid user",
			args:        args{"Hi, I'm a Valid [{'#`\"]}User Display Name"},
			wantValidID: true,
			want:        &User{displayname: "Hi, I'm a Valid [{'#`\"]}User Display Name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.displayname)
			if tt.wantValidID {
				tt.want.id = got.id
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRaw(t *testing.T) {
	type args struct {
		idStr       string
		displayname string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "should create a valid user",
			args: args{"123e4567-e89b-12d3-a456-426655440000", "user display name"},
			want: func() *User {
				id, _ := ParseID("123e4567-e89b-12d3-a456-426655440000")
				return &User{id, "user display name"}
			}(),
			wantErr: false,
		},
		{
			name:    "invalid uuid should return an error",
			args:    args{"invalid-uuid-format", "user display name"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRaw(tt.args.idStr, tt.args.displayname)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_ID(t *testing.T) {
	t.Run("should return a valid non-empty ID", func(t *testing.T) {
		u := New("display name")
		id := u.ID()
		if (id == ID{}) {
			t.Errorf("User.ID() = %v, should not be empty", id)
		}
	})

	t.Run("should return the same ID when creating user", func(t *testing.T) {
		initialID, err := ParseID("123e4567-e89b-12d3-a456-426655440000")
		if err != nil {
			t.Fatalf("Error parsing ID: %v", err)
		}
		u, err := NewRaw("123e4567-e89b-12d3-a456-426655440000", "display name")
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
		id := u.ID()
		if id != initialID {
			t.Errorf("User.ID() = %v, should equal %v", id, initialID)
		}
	})
}

func TestUser_DisplayName(t *testing.T) {
	t.Run("should return the same display name when creating user", func(t *testing.T) {
		initialDN := "display name"
		u, err := NewRaw("123e4567-e89b-12d3-a456-426655440000", "display name")
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
		dn := u.DisplayName()
		if dn != initialDN {
			t.Errorf("User.ID() = '%v', should equal '%v'", dn, initialDN)
		}
	})
}

func TestUser_UpdateDisplayName(t *testing.T) {
	type args struct {
		displayname string
	}
	tests := []struct {
		name string
		u    *User
		args args
		want string
	}{
		{
			name: "should update display name",
			u:    New("old name"),
			args: args{"new name"},
			want: "new name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.u.UpdateDisplayName(tt.args.displayname)
			got := tt.u.DisplayName()
			if got != tt.want {
				t.Errorf("User.UpdateDisplayName() updated to '%v', want '%v'", got, tt.want)
			}
		})
	}
}
