package ptt

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
)

func Test_getNewUtmpEnt(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uinfo *ptttype.UserInfoRaw
	}
	tests := []struct {
		name           string
		args           args
		expectedUtmpID ptttype.UtmpID
		sleepTS        time.Duration
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{testGetNewUtmpEnt0},
			expectedUtmpID: 5,
		},
		{
			args:           args{testGetNewUtmpEnt0},
			sleepTS:        time.Duration(1),
			expectedUtmpID: 5,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			time.Sleep(tt.sleepTS * time.Second)
			gotUtmpID, err := getNewUtmpEnt(tt.args.uinfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNewUtmpEnt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUtmpID, tt.expectedUtmpID) {
				t.Errorf("getNewUtmpEnt() = %v, want %v", gotUtmpID, tt.expectedUtmpID)
			}

			got := &cache.Shm.Shm.UInfo[gotUtmpID]

			testutil.TDeepEqual(t, "got", got, tt.args.uinfo)
		})
		wg.Wait()
	}
}

func TestGetUser(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID4 := &ptttype.UserID_t{}
	copy(userID4[:], []byte("SYSOP3"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *ptttype.UserecRaw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				userID: userID4,
			},
			expectedUser: testUserecRaw4,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUser, err := GetUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.expectedUser) {
				t.Errorf("GetUser() = %v, want %v", gotUser, tt.expectedUser)
			}
		})
		wg.Wait()
	}
}

func TestGetUserLevel(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name              string
		args              args
		expectedUserLevel ptttype.PERM
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args:              args{userID: &testUserecRaw1.UserID},
			expectedUserLevel: 536871943,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUserLevel, err := GetUserLevel(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUserLevel, tt.expectedUserLevel) {
				t.Errorf("GetUserLevel() = %v, want %v", gotUserLevel, tt.expectedUserLevel)
			}
		})
		wg.Wait()
	}
}

func TestGetUser2(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	testUserec2Raw1 := &ptttype.Userec2Raw{}

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name         string
		args         args
		expectedUser *ptttype.Userec2Raw
		wantErr      bool
	}{
		// TODO: Add test cases.
		{
			args:         args{userID: &testUserecRaw1.UserID},
			expectedUser: testUserec2Raw1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotUser, err := GetUser2(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.expectedUser) {
				t.Errorf("GetUser2() = %v, want %v", gotUser, tt.expectedUser)
			}
		})
		wg.Wait()
	}
}

func Test_postpermMsg(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	type args struct {
		uid   ptttype.UID
		user  *ptttype.UserecRaw
		bid   ptttype.Bid
		board *ptttype.BoardHeaderRaw
	}
	tests := []struct {
		name           string
		isNewBanSystem bool
		args           args
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:    "no pttype.PERM_POST",
			args:    args{uid: 2, user: testUserecRaw2, bid: 10, board: testBoardHeaderRaw2},
			wantErr: true,
		},
		{
			name:           "banned",
			args:           args{uid: 2, user: testNewPostUser2, bid: 10, board: testBoardHeaderRaw2},
			wantErr:        true,
			isNewBanSystem: true,
		},
		{
			name: "no ban in old system",
			args: args{uid: 2, user: testNewPostUser2, bid: 10, board: testBoardHeaderRaw2},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			ptttype.USE_NEW_BAN_SYSTEM = tt.isNewBanSystem
			defer func() {
				ptttype.USE_NEW_BAN_SYSTEM = true
			}()

			if err := postpermMsg(tt.args.uid, tt.args.user, tt.args.bid, tt.args.board); (err != nil) != tt.wantErr {
				t.Errorf("postpermMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
