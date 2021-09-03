package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
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

func Test_pwcuLoginSave(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	// 8.8.8.8
	ip := &ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56}

	// setup test case 1
	uid1 := ptttype.UID(1)
	uinfo1 := &ptttype.UserInfoRaw{
		UID:       uid1,
		UserLevel: 65535,
	}

	// setup test case 2
	userID3 := &ptttype.UserID_t{'u', 's', 'e', 'r', '3'}
	passwd3 := []byte("testpasswd")
	fromHost3 := &ptttype.IPv4_t{'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't'}
	email3 := &ptttype.Email_t{'t', 'e', 'm', 'p', '@', 't', 'e', 'm', 'p', '.', 'c', 'o', 'm'}
	nickname3 := &ptttype.Nickname_t{'n', 'i', 'c', 'k', 'n', 'a', 'm', 'e', '3'}
	realname3 := &ptttype.RealName_t{'r', 'e', 'a', 'l', 'n', 'a', 'm', 'e', '3'}
	career3 := &ptttype.Career_t{}
	address3 := &ptttype.Address_t{}
	uid3, _, err := Register(userID3, passwd3, fromHost3, email3, true, true, nickname3, realname3, career3, address3, true)
	logrus.Infof("after Register: e: %v", err)
	uinfo3 := &ptttype.UserInfoRaw{
		UID:       uid3,
		UserLevel: 7,
	}
	// setup test case 4
	uid5 := ptttype.UID(5)
	uinfo5 := &ptttype.UserInfoRaw{
		UID:       uid5,
		UserLevel: 7,
	}

	type args struct {
		uid   ptttype.UID
		uinfo *ptttype.UserInfoRaw
		ip    *ptttype.IPv4_t
	}
	tests := []struct {
		name                  string
		args                  args
		wantIsFirstLoginOfDay bool
		wantErr               bool
		wantNumLoginDay       uint32
		wantLastHost          ptttype.IPv4_t
	}{
		{
			name:                  "test user already login and first login this day, numLoginDays should increase 1, first login is true.",
			args:                  args{uid: uid1, uinfo: uinfo1, ip: ip},
			wantIsFirstLoginOfDay: true,
			wantErr:               false,
			wantNumLoginDay:       3,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user first login after register, numLoginDays should NOT increase, first login is false.",
			args:                  args{uid: uid3, uinfo: uinfo3, ip: ip},
			wantIsFirstLoginOfDay: false,
			wantErr:               false,
			wantNumLoginDay:       1,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user already login this day , numLoginDays should NOT increase, first login is false.",
			args:                  args{uid: uid1, uinfo: uinfo1, ip: ip},
			wantIsFirstLoginOfDay: false,
			wantErr:               false,
			wantNumLoginDay:       3,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user login after 1 day (23:59:59), numLoginDays should increase 1, first login is true.",
			args:                  args{uid: uid5, uinfo: uinfo5, ip: ip},
			wantIsFirstLoginOfDay: true,
			wantErr:               false,
			wantNumLoginDay:       2,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			userRaw, err := InitCurrentUserByUID(tt.args.uid)
			logrus.Infof("InitCurrentUserByUID: uid: %v userRaw: %v, e: %v", tt.args.uid, userRaw, err)
			gotIsFirstLoginOfDay, err := pwcuLoginSave(tt.args.uid, userRaw, tt.args.uinfo, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("pwcuLoginSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsFirstLoginOfDay != tt.wantIsFirstLoginOfDay {
				t.Errorf("pwcuLoginSave() gotIsFirstLoginOfDay = %v, want %v", gotIsFirstLoginOfDay, tt.wantIsFirstLoginOfDay)
			}

			if tt.wantNumLoginDay != userRaw.NumLoginDays {
				t.Errorf("pwcuLoginSave() gotNumLoginDay = %v, want %v", userRaw.NumLoginDays, tt.wantNumLoginDay)
			}

			testutil.TDeepEqual(t, "host", tt.wantLastHost, userRaw.LastHost)
		})
		wg.Wait()
	}
}
