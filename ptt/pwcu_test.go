package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func Test_pwcuRegCompleteJustify(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("CodingMan"))

	reg2 := &ptttype.Reg_t{}
	copy(reg2[:], "temp@temp.com")

	type args struct {
		uid     ptttype.UID
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
			args: args{ptttype.UID(2), userID2, reg2},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := pwcuRegCompleteJustify(tt.args.uid, tt.args.userID, tt.args.justify); (err != nil) != tt.wantErr {
				t.Errorf("pwcuRegCompleteJustify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func Test_pwcuBitDisableLevel(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], "CodingMan")

	type args struct {
		uid    ptttype.UID
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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := pwcuBitDisableLevel(tt.args.uid, tt.args.userID, tt.args.perm); (err != nil) != tt.wantErr {
				t.Errorf("pwcuBitDisableLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func Test_pwcuSetByBit(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

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

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if got := pwcuSetByBit(tt.args.perm, tt.args.mask, tt.args.isSet); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("pwcuSetByBit() = %v, want %v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}
