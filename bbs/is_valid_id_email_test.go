package bbs

import "testing"

func TestIsValidIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "A-pass",
			args: args{email: "abc@gmail.com"},
		},
		{
			name:    "A-fail",
			args:    args{email: "abd@gmail.com"},
			wantErr: true,
		},
		{
			name: "D-pass",
			args: args{email: "test@ptt.test"},
		},
		{
			name: "D-pass",
			args: args{email: "test@sub.ptt.test"},
		},
		{
			name: "P-pass",
			args: args{email: "test@cs.nthu.edu.tw"},
		},
		{
			name: "P-pass",
			args: args{email: "test@cs.nthu.edu.cn"},
		},
		{
			name: "S-pass",
			args: args{email: "test3@ntu.edu.tw"},
		},
		{
			name:    "S-fail",
			args:    args{email: "test@ntu.edu.sg"},
			wantErr: true,
		},
		{
			name:    "S-fail",
			args:    args{email: "test@csie.ntu.edu.tw"},
			wantErr: true,
		},
		{
			name:    "ban-P",
			args:    args{email: "test2@ptt.test"},
			wantErr: true,
		},
		{
			name:    "ban-A",
			args:    args{email: "test@ntu.edu.tw"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := IsValidIDEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("IsValidIDEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
