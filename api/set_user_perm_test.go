package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestSetUserPerm(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	origPerm := testUserec.Userlevel
	newPerm := ptttype.PERM_DEFAULT | ptttype.PERM_ADMIN | ptttype.PERM_LOGINOK

	params0 := &SetUserPermParams{Perm: newPerm}
	path0 := &SetUserPermPath{UserID: "SYSOP"}
	expected0 := &SetUserPermResult{Perm: newPerm}

	params1 := &SetUserPermParams{Perm: origPerm}
	expected1 := &SetUserPermResult{Perm: origPerm}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name           string
		args           args
		expectedResult interface{}
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params0, path: path0},
			expectedResult: expected0,
		},
		{
			args:           args{remoteAddr: testIP, uuserID: "SYSOP", params: params1, path: path0},
			expectedResult: expected1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := SetUserPerm(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUserPerm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("SetUserPerm() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
