package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

func TestSetUserPerm(t *testing.T) {
	setupTest()
	defer teardownTest()

	origPerm := testUserec1.Userlevel
	newPerm := ptttype.PERM_DEFAULT | ptttype.PERM_ADMIN | ptttype.PERM_LOGINOK

	logrus.Infof("TestSetUserPerm: origPerm: %v newPerm: %v", origPerm, newPerm)

	type args struct {
		userID    UUserID
		setUserID UUserID
		perm      ptttype.PERM
	}
	tests := []struct {
		name            string
		args            args
		expectedNewPerm ptttype.PERM
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{userID: "SYSOP", setUserID: "SYSOP", perm: newPerm},
			expectedNewPerm: newPerm,
		},
		{
			args:            args{userID: "SYSOP", setUserID: "SYSOP", perm: origPerm},
			expectedNewPerm: origPerm,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotNewPerm, err := SetUserPerm(tt.args.userID, tt.args.setUserID, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUserPerm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNewPerm, tt.expectedNewPerm) {
				t.Errorf("SetUserPerm() = %v, want %v", gotNewPerm, tt.expectedNewPerm)
			}
		})
		wg.Wait()
	}
}
