package cmbbs

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/sirupsen/logrus"
)

func TestPasswdLoadUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID1 := ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name      string
		args      args
		expected  ptttype.Uid
		expected1 *ptttype.UserecRaw
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			args:      args{&userID1},
			expected:  1,
			expected1: testUserecRaw1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			got, got1, err := PasswdLoadUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswdLoadUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("PasswdLoadUser() got = %v, expected %v", got, tt.expected)
			}
			testutil.TDeepEqual(t, "userec", got1, tt.expected1)
		})
		wg.Wait()
	}
}

func TestPasswdQuery(t *testing.T) {
	setupTest()
	defer teardownTest()
	type args struct {
		uid ptttype.Uid
	}
	tests := []struct {
		name     string
		args     args
		expected *ptttype.UserecRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{1},
			expected: testUserecRaw1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			got, err := PasswdQuery(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswdQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("PasswdQuery() = %v, expected %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestCheckPasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

	input1 := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', 0}
	input2 := []byte{'0', '1', '2', '4', '4', '5', '6', '7', '8', '9', '0', '1', 0}
	expected1 := ptttype.Passwd_t{65, 65, 51, 81, 66, 104, 76, 87, 107, 49, 66, 87, 65}

	type args struct {
		expected []byte
		input    []byte
	}
	tests := []struct {
		name     string
		args     args
		expected bool
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{expected1[:], input1},
			expected: true,
		},
		{
			name:     "incorrect input",
			args:     args{expected1[:], input2},
			expected: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			got, err := CheckPasswd(tt.args.expected, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("CheckPasswd() = %v, expected %v", got, tt.expected)
			}
		})
	}
	wg.Wait()
}

func TestGenPasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		passwd []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{[]byte("1234")},
		},
		{
			args: args{[]byte("abcdef")},
		},
		{
			args: args{[]byte("834792134")},
		},
		{
			args: args{[]byte("rweqrrwe")},
		},
		{
			args: args{[]byte("!@#$5ks")},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotPasswdHash, err := GenPasswd(tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenPasswd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			isGood, err := CheckPasswd(gotPasswdHash[:], tt.args.passwd)
			if err != nil || !isGood {
				t.Errorf("GenPasswd: unable to pass CheckPasswd: passwd: %v gotPasswdHash: %v", tt.args.passwd, gotPasswdHash)
			}
		})
		wg.Wait()
	}
}

func TestPasswdLock(t *testing.T) {
	setupTest()
	defer teardownTest()

	PasswdInit()
	defer PasswdDestroy()

	tests := []struct {
		name      string
		initSleep time.Duration
		sleep     time.Duration
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name:  "passwd-lock0",
			sleep: 5,
		},
		{
			name:      "passwd-lock1",
			initSleep: 6,
			sleep:     3,
		},
		{
			name:  "passwd-lock2",
			sleep: 5,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			nameChan := make(chan string)
			go func() {
				defer wg.Done()
				name := <-nameChan
				logrus.Infof("TestPasswdLock: %v: start", name)
				time.Sleep(tt.initSleep * time.Millisecond)
				if err := PasswdLock(); (err != nil) != tt.wantErr {
					t.Errorf("TestPasswdLock: PasswdLock() error = %v, wantErr %v", err, tt.wantErr)
				}
				logrus.Infof("TestPasswdLock: %v: got lock", name)
				time.Sleep(tt.sleep * time.Millisecond)
				logrus.Infof("TestPasswdLock: %v: done. to unlock", name)
				PasswdUnlock()
			}()
			nameChan <- tt.name
		})
	}
	wg.Wait()

}

func TestPasswdInit(t *testing.T) {
	setupTest()
	defer teardownTest()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PasswdInit(); (err != nil) != tt.wantErr {
				t.Errorf("PasswdInit() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer PasswdDestroy()
		})
		wg.Wait()
	}

}

func TestPasswdUpdate(t *testing.T) {
	setupTest()
	defer teardownTest()

	defer PasswdUpdate(1, testUserecRaw1)

	type args struct {
		uid  ptttype.Uid
		user *ptttype.UserecRaw
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected *ptttype.UserecRaw
	}{
		// TODO: Add test cases.
		{
			args:     args{1, testUserecRaw1Updated},
			expected: testUserecRaw1Updated,
		},
		{
			args:     args{1, testUserecRaw1},
			expected: testUserecRaw1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PasswdUpdate(tt.args.uid, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("PasswdUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, userecRaw, _ := PasswdLoadUser(&testUserecRaw1Updated.UserID)
			testutil.TDeepEqual(t, "userec", userecRaw, tt.expected)

		})
		wg.Wait()
	}
}

func TestPasswdUpdatePasswd(t *testing.T) {
	setupTest()
	defer teardownTest()

	passwd0 := []byte("123125")
	passwdHash0, _ := GenPasswd(passwd0)

	passwd1 := []byte("123123")
	passwdHash1, _ := GenPasswd(passwd1)

	type args struct {
		uid        ptttype.Uid
		passwdHash *ptttype.Passwd_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uid: 1, passwdHash: passwdHash0},
		},
		{
			args: args{uid: 1, passwdHash: passwdHash1},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PasswdUpdatePasswd(tt.args.uid, tt.args.passwdHash); (err != nil) != tt.wantErr {
				t.Errorf("PasswdUpdatePasswd() error = %v, wantErr %v", err, tt.wantErr)
			}

			passwdHash, _ := PasswdQueryPasswd(tt.args.uid)
			testutil.TDeepEqual(t, "passwdHash", passwdHash, tt.args.passwdHash)
		})
		wg.Wait()
	}
}

func TestPasswdUpdateEmail(t *testing.T) {
	setupTest()
	defer teardownTest()

	email0 := &ptttype.Email_t{}
	copy(email0[:], []byte("test@ptt.test"))

	type args struct {
		uid   ptttype.Uid
		email *ptttype.Email_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{uid: 1, email: email0},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PasswdUpdateEmail(tt.args.uid, tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("PasswdUpdateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}

			user, _ := PasswdQuery(tt.args.uid)

			testutil.TDeepEqual(t, "email", &user.Email, tt.args.email)
		})
	}
	wg.Wait()
}

func TestPasswdQueryUserLevel(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		uid ptttype.Uid
	}
	tests := []struct {
		name              string
		args              args
		expectedUserLevel ptttype.PERM
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{uid: 1},
			expectedUserLevel: 536871943,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserLevel, err := PasswdQueryUserLevel(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswdQueryUserLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserLevel, tt.expectedUserLevel) {
				t.Errorf("PasswdQueryUserLevel() = %v, want %v", gotUserLevel, tt.expectedUserLevel)
			}
		})
	}
	wg.Wait()
}

func TestPasswdUpdateUserLevel2(t *testing.T) {
	setupTest()
	defer teardownTest()

	userID := &ptttype.UserID_t{}
	copy(userID[:], []byte("CodingMan"))

	type args struct {
		userID *ptttype.UserID_t
		perm   ptttype.PERM2
		isSet  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{userID: userID, perm: ptttype.PERM2_ID_EMAIL, isSet: true},
		},
		{
			args: args{userID: userID, perm: ptttype.PERM2_ID_EMAIL, isSet: false},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := PasswdUpdateUserLevel2(tt.args.userID, tt.args.perm, tt.args.isSet); (err != nil) != tt.wantErr {
				t.Errorf("PasswdUpdateUserLevel2() error = %v, wantErr %v", err, tt.wantErr)
			}

			userLevel2, err := PasswdGetUserLevel2(tt.args.userID)
			if err != nil {
				t.Errorf("PasswdUpdateUserLevel2: e: %v", err)
			}
			if (tt.args.isSet && userLevel2&tt.args.perm == 0) || (!tt.args.isSet && userLevel2&tt.args.perm != 0) {
				t.Errorf("PasswdUpdateUserLevel2: isSet: %v userLevel2: (%v/%v/%v)", tt.args.isSet, userLevel2, tt.args.perm, userLevel2&tt.args.perm)
			}

			user2, err := PasswdGetUser2(tt.args.userID)
			if err != nil {
				t.Errorf("PasswdUpdateUserLevel2: e: %v", err)
			}

			if (tt.args.isSet && user2.UserLevel2&tt.args.perm == 0) || (!tt.args.isSet && user2.UserLevel2&tt.args.perm != 0) {
				t.Errorf("PasswdUpdateUserLevel2: isSet: %v userLevel2: (%v/%v/%v)", tt.args.isSet, userLevel2, tt.args.perm, userLevel2&tt.args.perm)
			}
		})
		wg.Wait()
	}
}

func TestPasswdGetUserLevel2(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name               string
		args               args
		expectedUserLevel2 ptttype.PERM2
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			args: args{&testUserecRaw1.UserID},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserLevel2, err := PasswdGetUserLevel2(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswdGetUserLevel2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserLevel2, tt.expectedUserLevel2) {
				t.Errorf("PasswdGetUserLevel2() = %v, want %v", gotUserLevel2, tt.expectedUserLevel2)
			}
		})
	}
	wg.Wait()
}
