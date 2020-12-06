package cmbbs

import (
	"reflect"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := PasswdLoadUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswdLoadUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("PasswdLoadUser() got = %v, expected %v", got, tt.expected)
			}
			if !reflect.DeepEqual(got1, tt.expected1) {
				t.Errorf("PasswdLoadUser() got1 = %v, expected1 %v", got1, tt.expected1)
			}
		})
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	wait := make(chan bool, len(tests))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameChan := make(chan string)
			go func() {
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
				wait <- true
			}()
			nameChan <- tt.name
		})
	}
	for i := 0; i < len(tests); i++ {
		<-wait
	}
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PasswdInit(); (err != nil) != tt.wantErr {
				t.Errorf("PasswdInit() error = %v, wantErr %v", err, tt.wantErr)
			}
			defer PasswdDestroy()
		})
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PasswdUpdate(tt.args.uid, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("PasswdUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, userecRaw, _ := PasswdLoadUser(&testUserecRaw1Updated.UserID)
			if !reflect.DeepEqual(userecRaw, tt.expected) {
				t.Errorf("UserecRaw: %v want: %v", userecRaw, tt.expected)
			}

		})
	}
}
