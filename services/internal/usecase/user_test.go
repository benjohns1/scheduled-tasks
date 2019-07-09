package usecase_test

import (
	"testing"

	data "github.com/benjohns1/scheduled-tasks/services/internal/data/transient"
	. "github.com/benjohns1/scheduled-tasks/services/internal/usecase"
)

func TestAddOrUpdateExternalUser(t *testing.T) {
	r := data.NewUserRepo()

	type args struct {
		r           UserRepo
		providerID  string
		externalID  string
		displayname string
	}
	type user struct {
		displayname string
	}
	tests := []struct {
		name    string
		args    args
		want    user
		wantErr ErrorCode
	}{
		{
			name:    "should add an empty user",
			args:    args{r, "provider", "externalID", ""},
			want:    user{""},
			wantErr: ErrNone,
		},
		{
			name:    "should add a basic user",
			args:    args{r, "provider", "externalID", "Test Displayname"},
			want:    user{"Test Displayname"},
			wantErr: ErrNone,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddOrUpdateExternalUser(tt.args.r, tt.args.providerID, tt.args.externalID, tt.args.displayname)
			if ((err == nil) != (tt.wantErr == ErrNone)) || ((err != nil) && (tt.wantErr != err.Code())) {
				t.Errorf("AddOrUpdateExternalUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.DisplayName() != tt.want.displayname {
				t.Errorf("AddOrUpdateExternalUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
