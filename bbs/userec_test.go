package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestNewUserecFromRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	testUserec2Raw := &ptttype.Userec2Raw{}

	type args struct {
		userecraw  *ptttype.UserecRaw
		userec2raw *ptttype.Userec2Raw
	}
	tests := []struct {
		name     string
		args     args
		expected *Userec
	}{
		// TODO: Add test cases.
		{
			args:     args{userecraw: testUserecRaw, userec2raw: testUserec2Raw},
			expected: testUserec1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got := NewUserecFromRaw(tt.args.userecraw, tt.args.userec2raw)

			testutil.TDeepEqual(t, "userec", got, tt.expected)
		})
	}
	wg.Wait()
}
