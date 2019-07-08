package user

import (
	"testing"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		name        string
		wantValidID bool
	}{
		{
			name:        "should return a valid ID",
			wantValidID: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewID()
			if (got != ID{}) != tt.wantValidID {
				t.Errorf("NewID() = %v, wantValidID %v", got, tt.wantValidID)
			}
		})
	}
}

func TestParseID(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name        string
		args        args
		wantValidID bool
		wantErr     bool
	}{
		{
			name:        "valid uuid string should return a valid ID",
			args:        args{"123e4567-e89b-12d3-a456-426655440000"},
			wantValidID: true,
			wantErr:     false,
		},
		{
			name:        "invalid uuid string should return an error",
			args:        args{"invalid-uuid-format"},
			wantValidID: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseID(tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != ID{}) != tt.wantValidID {
				t.Errorf("ParseID() = %v, wantValidID %v", got, tt.wantValidID)
			}
		})
	}
}

func parse(t *testing.T, val string) ID {
	id, err := ParseID(val)
	if err != nil {
		t.Fatalf("Error parsing ID val %v: %v", val, err)
	}
	return id
}

func TestID_Equals(t *testing.T) {
	id1 := NewID()
	id2 := NewID()
	id3a := parse(t, "123e4567-e89b-12d3-a456-426655440000")
	id3b := parse(t, "123e4567-e89b-12d3-a456-426655440000")
	id4 := parse(t, "111e1111-e89b-12d3-a456-426655440000")

	type args struct {
		other interface{}
	}
	tests := []struct {
		name string
		val  ID
		args args
		want bool
	}{
		{
			name: "the same new ID should be equal",
			val:  id1,
			args: args{id1},
			want: true,
		},
		{
			name: "zero value IDs should be equal",
			val:  ID{},
			args: args{ID{}},
			want: true,
		},
		{
			name: "different new IDs should not be equal",
			val:  id1,
			args: args{id2},
			want: false,
		},
		{
			name: "the same parsed IDs should be equal",
			val:  id3a,
			args: args{id3a},
			want: true,
		},
		{
			name: "IDs parsed from the same value should be equal",
			val:  id3a,
			args: args{id3b},
			want: true,
		},
		{
			name: "IDs parsed from different values should not be equal",
			val:  id3a,
			args: args{id4},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.Equals(tt.args.other); got != tt.want {
				t.Errorf("ID.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}
