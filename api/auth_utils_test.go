package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestVerifyJwt(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

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
		{ // XXX expires at 2031-07-14 06:33:12
			args:           args{raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiIiLCJleHAiOjE5NDE3OTE1OTIsInN1YiI6IlNZU09QIn0.hOUS-mN2KZeykTl8CMcJoFUafEC6o6LxNz7awisjhPI"},
			expectedUserID: bbs.UUserID("SYSOP"),
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserID, _, _, err := VerifyJwt(tt.args.raw, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserID, tt.expectedUserID) {
				t.Errorf("VerifyJwt() = %v, want %v", gotUserID, tt.expectedUserID)
			}
		})
		wg.Wait()
	}
}

func TestVerifyEmailJwt(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	// token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiIiLCJjdHgiOiJlbWFpbCIsImVtbCI6InRlc3RAcHR0LnRlc3QiLCJleHAiOjE5NDE3OTE4MDUsInN1YiI6IlNZU09QIn0.u5POkY8IrdzgQuEPnwaZM4gjadXW586c5uYAFmgISns
	token, _ := CreateEmailToken("SYSOP", "", "test@ptt.test", "email")

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
			gotUserID, _, gotClientInfo, gotEmail, err := VerifyEmailJwt(tt.args.raw, CONTEXT_CHANGE_EMAIL)
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
		wg.Wait()
	}
}

func Test_parseJwtClaim(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name       string
		args       args
		expectedCl *JwtClaim
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			args:       args{raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGkiOiIiLCJleHAiOjE5NDE3OTE1OTIsInN1YiI6IlNZU09QIn0.hOUS-mN2KZeykTl8CMcJoFUafEC6o6LxNz7awisjhPI"},
			expectedCl: &JwtClaim{UUserID: "SYSOP", Expire: 1941791592},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCl, err := parseJwtClaim(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJwtClaim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", gotCl, tt.expectedCl)
		})
	}
}
