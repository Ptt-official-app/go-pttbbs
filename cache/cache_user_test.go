package cache

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDoSearchUser(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestDoSearchUser: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	_ = LoadUHash()

	type args struct {
		userID   string
		isReturn bool
	}
	tests := []struct {
		name     string
		args     args
		expected ptttype.Uid
		want1    string
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			args:     args{userID: "SYSOP"},
			expected: 1,
			want1:    "",
		},
		{
			args:     args{userID: "SYSOP", isReturn: true},
			expected: 1,
			want1:    "SYSOP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := DoSearchUser(tt.args.userID, tt.args.isReturn)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSearchUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("DoSearchUser() got = %v, expected%v", got, tt.expected)
			}
			if got1 != tt.want1 {
				t.Errorf("DoSearchUser() got1 = %v, expected%v", got1, tt.want1)
			}
		})
	}
}

func TestAddToUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestDoSearchUser: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	InitFillUHash(false)

	user1 := &ptttype.UserID_t{}
	copy(user1[:], []byte("SYSOP"))

	user2 := &ptttype.UserID_t{}
	copy(user2[:], []byte("test1"))

	user3 := &ptttype.UserID_t{}
	copy(user3[:], []byte("test3"))

	type args struct {
		uidInCache ptttype.UidInStore
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddToUHash(tt.args.uidInCache, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("AddToUHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			userID, err := GetUserID(tt.args.uidInCache.ToUid())
			if err != nil {
				t.Errorf("AddToUHash: unable get user id: e: %v", err)
			}

			if !reflect.DeepEqual(userID, tt.args.userID) {
				t.Errorf("AddToUHash: userID: %v want: %v", userID, tt.args.userID)
			}

		})
	}
}

func TestRemoveFromUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestDoSearchUser: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	InitFillUHash(false)

	user := &ptttype.UserID_t{}

	AddToUHash(0, user)
	AddToUHash(1, user)
	AddToUHash(2, user)
	AddToUHash(3, user)
	AddToUHash(4, user)

	hashHead := &[1 << ptttype.HASH_BITS]ptttype.UidInStore{}
	nextInHash := &[ptttype.MAX_USERS]ptttype.UidInStore{}

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.HashHead),
		unsafe.Sizeof(Shm.Raw.HashHead),
		unsafe.Pointer(hashHead),
	)
	assert.Equal(t, ptttype.UidInStore(0), hashHead[35])

	Shm.ReadAt(
		unsafe.Offsetof(Shm.Raw.NextInHash),
		unsafe.Sizeof(Shm.Raw.NextInHash),
		unsafe.Pointer(nextInHash),
	)
	for i := 0; i < 4; i++ {
		assert.Equal(t, ptttype.UidInStore(i+1), nextInHash[i])
	}
	assert.Equal(t, ptttype.UidInStore(-1), nextInHash[4])
	for i := 5; i < len(nextInHash); i++ {
		assert.Equal(t, ptttype.UidInStore(0), nextInHash[i])
	}

	nextInHash1 := &[ptttype.MAX_USERS]ptttype.UidInStore{}
	copy(nextInHash1[:], []ptttype.UidInStore{1, 2, 3, 4, -1})

	nextInHash2 := &[ptttype.MAX_USERS]ptttype.UidInStore{}
	copy(nextInHash2[:], []ptttype.UidInStore{1, 3, 3, 4, -1})

	nextInHash3 := &[ptttype.MAX_USERS]ptttype.UidInStore{}
	copy(nextInHash3[:], []ptttype.UidInStore{1, 3, 3, -1, -1})

	type args struct {
		uidInHash ptttype.UidInStore
	}
	tests := []struct {
		name           string
		args           args
		wantHashHead   ptttype.UidInStore
		wantNextInHash *[ptttype.MAX_USERS]ptttype.UidInStore
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveFromUHash(tt.args.uidInHash); (err != nil) != tt.wantErr {
				t.Errorf("RemoveFromUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.HashHead),
				unsafe.Sizeof(Shm.Raw.HashHead),
				unsafe.Pointer(hashHead),
			)

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NextInHash),
				unsafe.Sizeof(Shm.Raw.NextInHash),
				unsafe.Pointer(nextInHash),
			)

			if !reflect.DeepEqual(hashHead[35], tt.wantHashHead) {
				t.Errorf("RemoveFromHash() hashHead: %v wantHashHead :%v", hashHead[35], tt.wantHashHead)
			}

			if !reflect.DeepEqual(nextInHash, tt.wantNextInHash) {
				t.Errorf("RemoveFromHash() nextInHash: %v wantNextInHash: %v", nextInHash, tt.wantNextInHash)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestGetUserID: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	InitFillUHash(false)

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))
	SetUserID(1, userID1)

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("SYSOP2"))
	SetUserID(2, userID2)

	userIDEmpty := &ptttype.UserID_t{}

	type args struct {
		uid ptttype.Uid
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserID(tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUserID(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("TestGetUserID: unable to NewSHM: e: %v", err)
		return
	}
	defer CloseSHM()

	InitFillUHash(false)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("SYSOP0"))

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("SYSOP"))

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("SYSOP2"))

	nextInHash1 := &[ptttype.MAX_USERS]int32{}
	copy(nextInHash1[:], []int32{-1})

	nextInHash2 := &[ptttype.MAX_USERS]int32{}
	copy(nextInHash2[:], []int32{-1, -1})

	type args struct {
		uid    ptttype.Uid
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name           string
		args           args
		wantUserID     *ptttype.UserID_t
		wantNextInHash *[ptttype.MAX_USERS]int32
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			nextInHash := &[ptttype.MAX_USERS]int32{}
			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NextInHash),
				unsafe.Sizeof(Shm.Raw.NextInHash),
				unsafe.Pointer(nextInHash),
			)
			assert.Equalf(t, nextInHash, tt.wantNextInHash, "SetUserID() nextInHash: %v want: %v", nextInHash, tt.wantNextInHash)
		})
	}
}
