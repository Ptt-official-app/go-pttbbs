package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptt/fav"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func Test_reginitFav(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	testFav0 := &fav.FavRaw{
		FavNum:  1,
		NBoards: 1,
		Favh: []*fav.FavType{
			{TheType: fav.FAVT_BOARD, Attr: 1, Fp: &fav.FavBoard{Bid: 1, LastVisit: 0, Attr: 0}},
		},
	}

	type args struct {
		uid  ptttype.Uid
		user *ptttype.UserecRaw
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		expected *fav.FavRaw
	}{
		// TODO: Add test cases.
		{
			args:     args{uid: 1, user: testUserecRaw1},
			expected: testFav0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := reginitFav(tt.args.uid, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("reginitFav() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _ := fav.Load(&tt.args.user.UserID)
			if got != nil {
				got.CleanRoot()
			}
			got.MTime = 0

			testutil.TDeepEqual(t, "got", got, tt.expected)
		})
		wg.Wait()
	}
}

func TestGetFavorites(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	userID0 := &ptttype.UserID_t{}
	copy(userID0[:], []byte("CodingMan"))

	expectedContent := []byte{
		0x23, 0x0d, 0x03, 0x00, 0x02, 0x01, 0x01, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x03, 0x01, 0x01, 0x02,
		0x01, 0x01, 0xb7, 0x73, 0xaa, 0xba, 0xa5, 0xd8,
		0xbf, 0xfd, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x01, 0x01,
		0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x01,
		0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	type args struct {
		userID     *ptttype.UserID_t
		retrieveTS types.Time4
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{userID: &testUserecRaw1.UserID},
		},
		{
			args:            args{userID: userID0},
			expectedContent: expectedContent,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, _, err := GetFavorites(tt.args.userID, tt.args.retrieveTS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFavorites() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("GetFavorites() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
		})
		wg.Wait()
	}
}
