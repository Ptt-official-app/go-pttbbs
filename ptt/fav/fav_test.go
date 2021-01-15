package fav

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFavLoad(t *testing.T) {
	setupTest()
	defer teardownTest()

	types.CopyFile("./testcase/home1", "./testcase/home")
	defer os.RemoveAll("./testcase/home")

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("testUser"))

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("testNoExist"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name     string
		args     args
		expected *FavRaw
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			name:     "testUser",
			args:     args{userID: userID0},
			expected: testFav0,
		},
		{
			name:     "testNoExist",
			args:     args{userID: userID1},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavLoad() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			logrus.Infof("TestFavLoad: after load: got: %v tt.expected: %v", got, tt.expected)
			if tt.expected != nil && got != nil {
				tt.expected.MTime = got.MTime
			}

			//XXX remove Root to avoid infinite recursive
			if got != nil {
				got.CleanRoot()
			}

			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
	}
}

func TestFavSave(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = types.CopyFile("./testcase/home1", "./testcase/home")
	defer os.RemoveAll("./testcase/home")

	time.Sleep(11 * time.Second)

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("testUserWrt"))

	userID1 := &ptttype.UserID_t{}
	copy(userID1[:], []byte("testUserWr1"))

	testFav2 := &FavRaw{
		FavNum:   5,
		NBoards:  2,
		NLines:   2,
		LineID:   2,
		NFolders: 1,
		FolderID: 1,
		Favh: []*FavType{
			{FAVT_LINE, 0, &FavLine{1}},
			{FAVT_LINE, 0, &FavLine{2}},
			{FAVT_FOLDER, 1, &FavFolder{1, testTitle0, testSubFav0}},
			{FAVT_BOARD, 1, &FavBoard{9, 0, 0}},
			{FAVT_BOARD, 1, &FavBoard{8, 0, 0}},
		},
	}

	gotFav2 := &FavRaw{
		FavNum:   4,
		NBoards:  2,
		NFolders: 1,
		FolderID: 1,
		Favh: []*FavType{
			{FAVT_FOLDER, 1, &FavFolder{1, testTitle0, testSubFav0}},
			{FAVT_BOARD, 1, &FavBoard{9, 0, 0}},
			{FAVT_BOARD, 1, &FavBoard{8, 0, 0}},
		},
	}

	userID2 := &ptttype.UserID_t{}
	copy(userID2[:], []byte("testUser2"))

	type args struct {
		fav    *FavRaw
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		sleep    time.Duration
		expected *FavRaw
	}{
		// TODO: Add test cases.
		{
			name:     "test0",
			args:     args{fav: testFav0, userID: userID0},
			expected: testFav0,
		},
		{
			name:    "test1",
			args:    args{fav: testFav1, userID: userID1},
			wantErr: true,
		},
		{
			name:     "test2",
			args:     args{fav: testFav2, userID: userID2},
			expected: gotFav2,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			time.Sleep(tt.sleep)
			tt.args.fav.MTime = types.NowTS()

			got, err := tt.args.fav.Save(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavSave() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got == nil {
				return
			}

			//XXX remove Root to avoid infinite recursive
			if got != nil {
				got.CleanRoot()
			}

			tt.expected.MTime = got.MTime
			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestFavRaw_AddBoard(t *testing.T) {
	f0 := NewFavRaw(nil)
	ft0 := &FavType{FAVT_BOARD, 1, &FavBoard{1, 0, 0}}
	type args struct {
		bid ptttype.Bid
	}

	expected0 := &FavRaw{
		FavNum:  1,
		NBoards: 1,
		Favh: []*FavType{
			{FAVT_BOARD, 1, &FavBoard{Bid: 1}},
		},
	}

	tests := []struct {
		name            string
		f               *FavRaw
		args            args
		expectedFavType *FavType
		expected        *FavRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			f:               f0,
			args:            args{bid: 1},
			expectedFavType: ft0,
			expected:        expected0,
		},
		{
			f:               f0,
			args:            args{bid: 1},
			expectedFavType: ft0,
			expected:        expected0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			f := tt.f
			gotFavType, err := f.AddBoard(tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("FavRaw.AddBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFavType, tt.expectedFavType) {
				t.Errorf("FavRaw.AddBoard() = %v, want %v", gotFavType, tt.expectedFavType)
			}

			f.CleanRoot()

			testutil.TDeepEqual(t, "expected", f, tt.expected)

		})
		wg.Wait()
	}
}

func TestFavRaw_AddLine(t *testing.T) {
	f0 := NewFavRaw(nil)
	ft0 := &FavType{FAVT_LINE, 1, &FavLine{1}}

	expected0 := &FavRaw{
		FavNum: 1,
		NLines: 1,
		LineID: 1,
		Favh: []*FavType{
			{FAVT_LINE, 1, &FavLine{1}},
		},
	}

	tests := []struct {
		name            string
		f               *FavRaw
		expectedFavType *FavType
		expected        *FavRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			f:               f0,
			expectedFavType: ft0,
			expected:        expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.f
			gotFavType, err := f.AddLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("FavRaw.AddLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFavType, tt.expectedFavType) {
				t.Errorf("FavRaw.AddLine() = %v, want %v", gotFavType, tt.expectedFavType)
			}

			f.CleanRoot()

			testutil.TDeepEqual(t, "expected", f, tt.expected)
		})
	}
}

func TestFavRaw_AddFolder(t *testing.T) {
	f0 := NewFavRaw(nil)
	ft0 := &FavType{FAVT_FOLDER, 1, &FavFolder{Fid: 1, ThisFolder: NewFavRaw(nil)}}

	expected0 := &FavRaw{
		FavNum:   1,
		NFolders: 1,
		FolderID: 1,
		Favh: []*FavType{
			{FAVT_FOLDER, 1, &FavFolder{Fid: 1, ThisFolder: NewFavRaw(nil)}},
		},
	}

	tests := []struct {
		name            string
		f               *FavRaw
		expectedFavType *FavType
		expected        *FavRaw
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			f:               f0,
			expectedFavType: ft0,
			expected:        expected0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.f
			gotFavType, err := f.AddFolder()
			if (err != nil) != tt.wantErr {
				t.Errorf("FavRaw.AddFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			f.CleanRoot()
			tt.expected.CleanRoot()

			gotFolder := gotFavType.Fp.(*FavFolder)
			gotFolder.ThisFolder.CleanRoot()

			expectedFolder := tt.expectedFavType.Fp.(*FavFolder)
			expectedFolder.ThisFolder.CleanRoot()

			assert.Equal(t, gotFavType, tt.expectedFavType)

			testutil.TDeepEqual(t, "expected", f, tt.expected)
		})
	}
}
