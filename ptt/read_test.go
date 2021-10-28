package ptt

import (
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestDeleteArticles(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))

	filename0 := &ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1607202239.A.30D"))

	type args struct {
		boardID  *ptttype.BoardID_t
		filename *ptttype.Filename_t
		index    ptttype.SortIdx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test DeleteArticles index 0 works w/o error",
			args: args{
				boardID0,
				filename0,
				1,
			},
			wantErr: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			wg.Done()
			if err := DeleteArticles(tt.args.boardID, tt.args.filename, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("DeleteArticles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
