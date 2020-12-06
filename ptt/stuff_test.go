package ptt

import (
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_is_uBM(t *testing.T) {
	bm := &ptttype.BM_t{}
	user0 := &ptttype.UserID_t{}
	user1 := &ptttype.UserID_t{}
	user2 := &ptttype.UserID_t{}
	user3 := &ptttype.UserID_t{}
	user4 := &ptttype.UserID_t{}

	copy(bm[:], "abcde/12345/B2S")
	copy(user0[:], "12345")
	copy(user1[:], "234")
	copy(user2[:], "e/12345")
	copy(user3[:], "abc")
	copy(user4[:], "b2s")

	type args struct {
		userID *ptttype.UserID_t
		bm     *ptttype.BM_t
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		// TODO: Add test cases.
		{
			args:     args{user0, bm},
			expected: true,
		},
		{
			args:     args{user1, bm},
			expected: false,
		},
		{
			args:     args{user2, bm},
			expected: false,
		},
		{
			args:     args{user3, bm},
			expected: false,
		},
		{
			name:     "require exact match",
			args:     args{user4, bm},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := is_uBM(tt.args.userID, tt.args.bm); got != tt.expected {
				t.Errorf("is_uBM() = %v, want %v", got, tt.expected)
			}
		})
	}
}
