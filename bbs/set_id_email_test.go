package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestSetIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uuserID UUserID
		isSet   bool
	}
	tests := []struct {
		name               string
		args               args
		expectedUserLevel2 ptttype.PERM2
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args:               args{uuserID: "SYSOP", isSet: true},
			expectedUserLevel2: ptttype.PERM2_ID_EMAIL,
		},
		{
			args:               args{uuserID: "SYSOP", isSet: false},
			expectedUserLevel2: ptttype.PERM2_INVALID,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserLevel2, err := SetIDEmail(tt.args.uuserID, tt.args.isSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIDEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserLevel2, tt.expectedUserLevel2) {
				t.Errorf("SetIDEmail() = %v, want %v", gotUserLevel2, tt.expectedUserLevel2)
			}
		})
		wg.Wait()
	}
}
