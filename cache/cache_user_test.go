package cache

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/stretchr/testify/assert"
)

func TestDoSearchUserRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = LoadUHash()

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name     string
		args     args
		expected ptttype.UID
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: userID0},
			expected: 1,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := DoSearchUserRaw(tt.args.userID, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("DoSearchUser() got = %v, expected%v", got, tt.expected)
			}
		})
		wg.Wait()
	}
}

func TestAddToUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	user1 := &ptttype.UserID_t{}
	copy(user1[:], []byte("SYSOP"))

	user2 := &ptttype.UserID_t{}
	copy(user2[:], []byte("test1"))

	user3 := &ptttype.UserID_t{}
	copy(user3[:], []byte("test3"))

	type args struct {
		uidInCache ptttype.UIDInStore
		userID     *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{0, user1},
		},
		{
			args: args{1, user2},
		},
		{
			args: args{2, user3},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := AddToUHash(tt.args.uidInCache, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("AddToUHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			userID, err := GetUserID(tt.args.uidInCache.ToUID())
			if err != nil {
				t.Errorf("AddToUHash: unable get user id: e: %v", err)
			}

			if !reflect.DeepEqual(userID, tt.args.userID) {
				t.Errorf("AddToUHash: userID: %v want: %v", userID, tt.args.userID)
			}
		})
		wg.Wait()
	}
}

func TestRemoveFromUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	user := &ptttype.UserID_t{}

	AddToUHash(0, user)
	AddToUHash(1, user)
	AddToUHash(2, user)
	AddToUHash(3, user)
	AddToUHash(4, user)

	hashHead := &Shm.Shm.HashHead
	assert.Equal(t, ptttype.UIDInStore(0), hashHead[35])

	nextInHash := &Shm.Shm.NextInHash
	for i := 0; i < 4; i++ {
		assert.Equal(t, ptttype.UIDInStore(i+1), nextInHash[i])
	}
	assert.Equal(t, ptttype.UIDInStore(-1), nextInHash[4])
	for i := 5; i < len(nextInHash); i++ {
		assert.Equal(t, ptttype.UIDInStore(0), nextInHash[i])
	}

	nextInHash1 := &[ptttype.MAX_USERS]ptttype.UIDInStore{}
	copy(nextInHash1[:], []ptttype.UIDInStore{1, 2, 3, 4, -1})

	nextInHash2 := &[ptttype.MAX_USERS]ptttype.UIDInStore{}
	copy(nextInHash2[:], []ptttype.UIDInStore{1, 3, 3, 4, -1})

	nextInHash3 := &[ptttype.MAX_USERS]ptttype.UIDInStore{}
	copy(nextInHash3[:], []ptttype.UIDInStore{1, 3, 3, -1, -1})

	type args struct {
		uidInHash ptttype.UIDInStore
	}
	tests := []struct {
		name           string
		args           args
		wantHashHead   ptttype.UIDInStore
		wantNextInHash *[ptttype.MAX_USERS]ptttype.UIDInStore
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			args:           args{0},
			wantHashHead:   1,
			wantNextInHash: nextInHash1,
		},
		{
			name:           "dupped removing cases",
			args:           args{0},
			wantHashHead:   1,
			wantNextInHash: nextInHash1,
		},
		{
			name:           "with 1, 3, 4",
			args:           args{2},
			wantHashHead:   1,
			wantNextInHash: nextInHash2,
		},
		{
			name:           "drop last",
			args:           args{4},
			wantHashHead:   1,
			wantNextInHash: nextInHash3,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := RemoveFromUHash(tt.args.uidInHash); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(hashHead[35], tt.wantHashHead) {
				t.Errorf("RemoveFromHash() hashHead: %v wantHashHead :%v", hashHead[35], tt.wantHashHead)
			}

			if !reflect.DeepEqual(nextInHash, tt.wantNextInHash) {
				t.Errorf("RemoveFromHash() nextInHash: %v wantNextInHash: %v", nextInHash, tt.wantNextInHash)
			}
		})
		wg.Wait()
	}
}

func TestGetUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))
	SetUserID(1, userID1)

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("SYSOP2"))
	SetUserID(2, userID2)

	userIDEmpty := &ptttype.UserID_t{}

	type args struct {
		uid ptttype.UID
	}
	tests := []struct {
		name    string
		args    args
		want    *ptttype.UserID_t
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args:    args{0},
			wantErr: true,
		},
		{
			args: args{1},
			want: userID1,
		},
		{
			args: args{2},
			want: userID2,
		},
		{
			args: args{3},
			want: userIDEmpty,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := GetUserID(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserID() = %v, want %v", got, tt.want)
			}
		})
		wg.Wait()
	}
}

func TestSetUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP0"))

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("SYSOP2"))

	nextInHash1 := &[ptttype.MAX_USERS]ptttype.UIDInStore{}
	copy(nextInHash1[:], []ptttype.UIDInStore{-1})

	nextInHash2 := &[ptttype.MAX_USERS]ptttype.UIDInStore{}
	copy(nextInHash2[:], []ptttype.UIDInStore{-1, -1})

	type args struct {
		uid    ptttype.UID
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name           string
		args           args
		wantUserID     *ptttype.UserID_t
		wantNextInHash *[ptttype.MAX_USERS]ptttype.UIDInStore
		wantErr        bool
	}{
		// TODO: Add test cases.
		{
			name:    "invalid userID",
			args:    args{0, userID0},
			wantErr: true,
		},
		{
			args:           args{1, userID1},
			wantNextInHash: nextInHash1,
		},
		{
			args:           args{1, userID1},
			wantNextInHash: nextInHash1,
		},
		{
			args:           args{2, userID2},
			wantNextInHash: nextInHash2,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			err := SetUserID(tt.args.uid, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUserID() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			userID, _ := GetUserID(tt.args.uid)
			if !reflect.DeepEqual(userID, tt.args.userID) {
				t.Errorf("SetUserID() userID: %v want: %v", userID, tt.args.userID)
			}

			nextInHash := &Shm.Shm.NextInHash
			assert.Equalf(t, nextInHash, tt.wantNextInHash, "SetUserID() nextInHash: %v want: %v", nextInHash, tt.wantNextInHash)
		})
		wg.Wait()
	}
}

func TestAddCooldownTime(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	type args struct {
		uid     ptttype.UID
		minutes int
	}
	tests := []struct {
		name     string
		args     args
		expected int
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{uid: 1, minutes: 10},
			expected: 10,
		},
		{
			args:     args{uid: 1, minutes: 5},
			expected: 15,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := AddCooldownTime(tt.args.uid, tt.args.minutes); (err != nil) != tt.wantErr {
				t.Errorf("AddCooldownTime() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := CooldownTimeOf(tt.args.uid)
			nowTS := types.NowTS()
			if !((nowTS+types.Time4(tt.expected-1)*60) <= got && got <= (nowTS+types.Time4(tt.expected*60))) {
				t.Errorf("AddCooldownTime: got: %v nowTS: %v diff: %v expected: %v", got, nowTS, got-nowTS, tt.args.minutes*60)
			}
		})
		wg.Wait()
	}
}

func TestAddPosttimes(t *testing.T) {
	setupTest()
	defer teardownTest()

	InitFillUHash(false)

	type args struct {
		uid   ptttype.UID
		times int
	}
	tests := []struct {
		name     string
		args     args
		expected types.Time4
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{uid: 1, times: 1},
			expected: 1,
		},
		{
			args:     args{uid: 1, times: 2},
			expected: 3,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := AddPosttimes(tt.args.uid, tt.args.times); (err != nil) != tt.wantErr {
				t.Errorf("AddPosttimes() error = %v, wantErr %v", err, tt.wantErr)
			}

			posttimes := PosttimesOf(tt.args.uid)
			if posttimes != tt.expected {
				t.Errorf("AddPosttimes: posttimes: %v expected: %v", posttimes, tt.expected)
			}
		})
		wg.Wait()
	}
}
