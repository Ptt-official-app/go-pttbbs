package ptt

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_pwcuRegCompleteJustify(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("CodingMan"))

	reg2 := &ptttype.Reg_t{}
	copy(reg2[:], "temp@temp.com")

	type args struct {
		uid     ptttype.Uid
		userID  *ptttype.UserID_t
		justify *ptttype.Reg_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{ptttype.Uid(2), userID2, reg2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pwcuRegCompleteJustify(tt.args.uid, tt.args.userID, tt.args.justify); (err != nil) != tt.wantErr {
				t.Errorf("pwcuRegCompleteJustify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_pwcuBitDisableLevel(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], "CodingMan")

	type args struct {
		uid    ptttype.Uid
		userID *ptttype.UserID_t
		perm   ptttype.PERM
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{2, userID2, ptttype.PERM_BASIC},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := pwcuBitDisableLevel(tt.args.uid, tt.args.userID, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("pwcuBitDisableLevel() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func Test_pwcuSetByBit(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		perm  ptttype.PERM
		mask  ptttype.PERM
		isSet bool
	}
	tests := []struct {
		name     string
		args     args
		expected ptttype.PERM
	}{
		// TODO: Add test cases.
		{
			args:     args{ptttype.PERM_BASIC, ptttype.PERM_POST, true},
			expected: ptttype.PERM_BASIC | ptttype.PERM_POST,
		},
		{
			args:     args{ptttype.PERM_BASIC | ptttype.PERM_POST, ptttype.PERM_POST, false},
			expected: ptttype.PERM_BASIC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pwcuSetByBit(tt.args.perm, tt.args.mask, tt.args.isSet); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("pwcuSetByBit() = %v, want %v", got, tt.expected)
			}
		})
	}
}
