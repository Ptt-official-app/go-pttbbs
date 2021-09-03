package bbs

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptt"
)

func TestReloadUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	_ = ptt.SetupNewUser(testPermissionUserecRaw)
	_ = ptt.SetupNewUser(testPermissionUserecRaw2)
	_ = ptt.SetupNewUser(testPermissionUserecRaw3)

	type args struct {
		userID UUserID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "For SYSOP, should be work",
			args:    args{userID: "test"},
			wantErr: false,
		},
		{
			name:    "For A1, should NOT be work (0)",
			args:    args{userID: "test1"},
			wantErr: true,
		},
		{
			name:    "For test, should NOT be work (ptttype.PERM_BOARD | ptttype.PERM_POST | ptttype.PERM_LOGINOK)",
			args:    args{userID: "test2"},
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
