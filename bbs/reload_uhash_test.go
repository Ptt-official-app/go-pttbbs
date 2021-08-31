package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestReloadUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testPermissionUserecRaw)
	_ = ptt.SetupNewUser(testPermissionUserecRaw2)
	_ = ptt.SetupNewUser(testPermissionUserecRaw3)

	// setup test case 1: For test, should be work
	userID1 := ptttype.UserID_t{}
	copy(userID1[:], []byte("test"))
	// setup test case 2: For test1, should NOT be work
	userID2 := ptttype.UserID_t{}
	copy(userID2[:], []byte("test1"))
	// setup test case 3: For test2, should NOT be work
	userID3 := ptttype.UserID_t{}
	copy(userID3[:], []byte("test2"))

	type args struct {
		userID *ptttype.UserID_t
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "For SYSOP, should be work",
			args:    args{userID: &userID1},
			wantErr: false,
		},
		{
			name:    "For A1, should NOT be work (0)",
			args:    args{userID: &userID2},
			wantErr: true,
		},
		{
			name:    "For test, should NOT be work (ptttype.PERM_BOARD | ptttype.PERM_POST | ptttype.PERM_LOGINOK)",
			args:    args{userID: &userID3},
			wantErr: true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := ReloadUHash(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("ReloadUHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
