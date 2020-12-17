package bbs

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestNewUserecFromRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uid       ptttype.Uid
		userecraw *ptttype.UserecRaw
	}
	tests := []struct {
		name     string
		args     args
		expected *Userec
	}{
		// TODO: Add test cases.
		{
			args:     args{uid: 1, userecraw: testUserecRaw},
			expected: testUserec1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserecFromRaw(tt.args.uid, tt.args.userecraw)

			types.TDeepEqual(t, got, tt.expected)
		})
	}
}
