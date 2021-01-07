package bbs

import "testing"

func TestChangeEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		userID UUserID
		email  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{userID: "SYSOP", email: "test@ptt.test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeEmail(tt.args.userID, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("ChangeEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, _ := GetUser(tt.args.userID)

			if user.Email != tt.args.email {
				t.Errorf("ChangeEmail: %v email: %v", user.Email, tt.args.email)
			}
		})
	}
}
