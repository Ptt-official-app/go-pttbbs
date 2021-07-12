package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestVerifyJwt(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		raw string
	}
	tests := []struct {
		name           string
		args           args
		expectedUserID bbs.UUserID
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{},
			expectedUserID: bbs.UUserID(GUEST),
		},
		{
			args:    args{raw: "not-exists"},
			wantErr: true,
		},
		{
			args:    args{raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFeHBpcmUiOjE2MDgzMzA0NTYsIlVzZXJJRCI6IlNZU09QMiJ9.G6gKhrGRysMAvOJb6rMmsvqxm7MuUwOkHhII7D73Ijc"},
			wantErr: true,
		},
		{ //XXX expires at 2270-11-01 09:53:32
			args:           args{raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiIiLCJleHAiOjk0OTM0MjI4MTIsInN1YiI6IlNZU09QMiJ9.VbixNBxg4h5FCyTmvhtVzBJ4HsE5_va-MPZzR8TLaY8"},
			expectedUserID: bbs.UUserID("SYSOP2"),
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserID, _, err := VerifyJwt(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserID, tt.expectedUserID) {
				t.Errorf("VerifyJwt() = %v, want %v", gotUserID, tt.expectedUserID)
			}
		})
	}
	wg.Wait()
}

func TestVerifyEmailJwt(t *testing.T) {
	setupTest()
	defer teardownTest()

	token, _ := CreateEmailToken("SYSOP", "", "test@ptt.test", CONTEXT_CHANGE_EMAIL)

	type args struct {
		raw string
	}
	tests := []struct {
		name               string
		args               args
		expectedUserID     bbs.UUserID
		expectedClientInfo string
		expectedEmail      string
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args:           args{raw: token},
			expectedUserID: "SYSOP",
			expectedEmail:  "test@ptt.test",
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserID, gotClientInfo, gotEmail, err := VerifyEmailJwt(tt.args.raw, CONTEXT_CHANGE_EMAIL)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyEmailJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserID, tt.expectedUserID) {
				t.Errorf("VerifyEmailJwt() gotUserID = %v, want %v", gotUserID, tt.expectedUserID)
			}
			if gotClientInfo != tt.expectedClientInfo {
				t.Errorf("VerifyEmailJwt() gotClientInfo = %v, want %v", gotClientInfo, tt.expectedClientInfo)
			}
			if gotEmail != tt.expectedEmail {
				t.Errorf("VerifyEmailJwt() gotEmail = %v, want %v", gotEmail, tt.expectedEmail)
			}
		})
	}
	wg.Wait()
}
