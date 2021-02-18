package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestSetIDEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	jwt, _ := CreateEmailToken("SYSOP", "", "test@ptt.test", CONTEXT_SET_ID_EMAIL)

	params0 := &SetIDEmailParams{
		IsSet: true,
		Jwt:   jwt,
	}

	path0 := &SetIDEmailPath{
		UserID: "SYSOP",
	}

	result0 := &SetIDEmailResult{
		UserID:     "SYSOP",
		Email:      "test@ptt.test",
		UserLevel2: ptttype.PERM2_ID_EMAIL,
	}

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
			expectedResult: result0,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotResult, err := SetIDEmail(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIDEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.expectedResult) {
				t.Errorf("SetIDEmail() = %v, want %v", gotResult, tt.expectedResult)
			}
		})
		wg.Wait()
	}
}
