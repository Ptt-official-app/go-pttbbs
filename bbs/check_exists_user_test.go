package bbs

import "testing"

func TestCheckExistsUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		username string
	}
	tests := []struct {
		name            string
		args            args
		expectedUuserID UUserID
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{username: "SYSOP"},
			expectedUuserID: "SYSOP",
		},
		{
			args: args{username: "notexists"},
		},
		{
			args:    args{username: "not-exists"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUuserID, err := CheckExistsUser(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExistsUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUuserID != tt.expectedUuserID {
				t.Errorf("CheckExistsUser() = %v, want %v", gotUuserID, tt.expectedUuserID)
			}
		})
	}
}
