package ptt

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/types"

	"github.com/sirupsen/logrus"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func TestInitCurrentUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userid1 := ptttype.UserID_t{}
	copy(userid1[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name      string
		args      args
		expected  ptttype.UID
		expected1 *ptttype.UserecRaw
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			args:      args{&userid1},
			expected:  1,
			expected1: testUserecRaw1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, got1, err := InitCurrentUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitCurrentUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("InitCurrentUser() got = %v, expected %v", got, tt.expected)
			}
			if !reflect.DeepEqual(got1, tt.expected1) {
				t.Errorf("InitCurrentUser() got1 = %v, expected1 %v", got1, tt.expected1)
			}
		})
		wg.Wait()
	}
}

func Test_passwdSyncUpdate(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uid  ptttype.UID
		user *ptttype.UserecRaw
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uid: ptttype.UID(1), user: testUserecRaw1},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := passwdSyncUpdate(tt.args.uid, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("passwdSyncUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func Test_passwdSyncQuery(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uid ptttype.UID
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.UserecRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{11},
			expected: testUserecRaw4,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := passwdSyncQuery(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("passwdSyncQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestInitCurrentUserByUid(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uid ptttype.UID
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *ptttype.UserecRaw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{1},
			expectedUser: testUserecRaw1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUser, err := InitCurrentUserByUID(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitCurrentUserByUid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.expectedUser) {
				t.Errorf("InitCurrentUserByUid() = %v, want %v", gotUser, tt.expectedUser)
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
	userID1 := ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))
	uID1, userRaw1, _ := InitCurrentUser(&userID1)
	logrus.Infof("firstLogin: %v lastLogin: %v numLoginDays: %v\n",
		userRaw1.FirstLogin, userRaw1.LastLogin, userRaw1.NumLoginDays)
	// setup test case 2
	uID3 := ptttype.UID(3)
	userRaw3, _ := InitCurrentUserByUID(uID3)
	userRaw3.NumLoginDays = 0
	userRaw3.FirstLogin = 1630000000 // Thu Aug 26 2021 17:46:40 GMT+0000
	userRaw3.LastLogin = 0
	// setup test case 3
	now := time.Now()
	baseNowDay := types.TimeToTime4(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))
	uID4 := ptttype.UID(4)
	userRaw4, _ := InitCurrentUserByUID(uID4)
	userRaw4.NumLoginDays = 1
	userRaw4.FirstLogin = 1630000000 // Thu Aug 26 2021 17:46:40 GMT+0000
	userRaw4.LastLogin = baseNowDay + 600
	// setup test case 4
	uID5 := ptttype.UID(5)
	userRaw5, _ := InitCurrentUserByUID(uID5)
	userRaw5.NumLoginDays = 1
	userRaw5.FirstLogin = 1630000000 // Thu Aug 26 2021 17:46:40 GMT+0000
	userRaw5.LastLogin = baseNowDay - 1

	type args struct {
		uid  ptttype.UID
		user *ptttype.UserecRaw
		ip   *ptttype.IPv4_t
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
			args:                  args{uid: uID1, user: userRaw1, ip: ip},
			wantIsFirstLoginOfDay: true,
			wantErr:               false,
			wantNumLoginDay:       3,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user first login after register, numLoginDays should increase 1, first login is true.",
			args:                  args{uid: uID3, user: userRaw3, ip: ip},
			wantIsFirstLoginOfDay: true,
			wantErr:               false,
			wantNumLoginDay:       1,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user already login this day , numLoginDays should NOT increase, first login is false.",
			args:                  args{uid: uID4, user: userRaw4, ip: ip},
			wantIsFirstLoginOfDay: false,
			wantErr:               false,
			wantNumLoginDay:       1,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
		{
			name:                  "test user login after 1 day (23:59:59), numLoginDays should increase 1, first login is true.",
			args:                  args{uid: uID5, user: userRaw5, ip: ip},
			wantIsFirstLoginOfDay: true,
			wantErr:               false,
			wantNumLoginDay:       2,
			wantLastHost:          ptttype.IPv4_t{56, 46, 56, 46, 56, 46, 56},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsFirstLoginOfDay, err := pwcuLoginSave(tt.args.uid, tt.args.user, tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("pwcuLoginSave() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsFirstLoginOfDay != tt.wantIsFirstLoginOfDay {
				t.Errorf("pwcuLoginSave() gotIsFirstLoginOfDay = %v, want %v", gotIsFirstLoginOfDay, tt.wantIsFirstLoginOfDay)
			}
			if tt.wantNumLoginDay != tt.args.user.NumLoginDays {
				t.Errorf("pwcuLoginSave() gotNumLoginDay = %v, want %v", tt.args.user.NumLoginDays, tt.wantNumLoginDay)
			}
			if tt.wantLastHost != tt.args.user.LastHost {
				t.Errorf("pwcuLoginSave() gotLastHost = %v, want %v", tt.args.user.LastHost, tt.wantLastHost)
			}
		})
	}
}
