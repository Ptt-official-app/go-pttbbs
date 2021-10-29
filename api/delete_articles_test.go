package api

import (
	"encoding/binary"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/types"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestDeleteArticles(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())
	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("10_WhoAmI"))

	filename0 := &ptttype.Filename_t{}
	copy(filename0[:], []byte("M.1607202239.A.30D"))

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607203395.A.00D"))

	fileHeaderRaw1 := &ptttype.FileHeaderRaw{
		Filename: ptttype.Filename_t{ // M.1607202239.A.30D
			0x4d, 0x2e, 0x31, 0x36, 0x30, 0x37, 0x32, 0x30,
			0x32, 0x32, 0x33, 0x39, 0x2e, 0x41, 0x2e, 0x33,
			0x30, 0x44,
		},
		Modified: 1607202238,
		Owner:    ptttype.Owner_t{0x53, 0x59, 0x53, 0x4f, 0x50}, // SYSOP
		Date:     ptttype.Date_t{0x31, 0x32, 0x2f, 0x30, 0x36},  // 12/06
		Title: ptttype.Title_t{ //[問題] 我是誰？～
			0x5b, 0xb0, 0xdd, 0xc3, 0x44, 0x5d, 0x20, 0xa7,
			0xda, 0xac, 0x4f, 0xbd, 0xd6, 0xa1, 0x48, 0xa1,
			0xe3,
		},
	}
	fileHeaderRaw2 := &ptttype.FileHeaderRaw{
		Filename: ptttype.Filename_t{ // M.1607203395.A.00D
			0x4d, 0x2e, 0x31, 0x36, 0x30, 0x37, 0x32, 0x30,
			0x33, 0x33, 0x39, 0x35, 0x2e, 0x41, 0x2e, 0x30,
			0x30, 0x44,
		},
		Modified: 1607203394,
		Owner:    ptttype.Owner_t{0x53, 0x59, 0x53, 0x4f, 0x50}, // SYSOP
		Date:     ptttype.Date_t{0x31, 0x32, 0x2f, 0x30, 0x36},  // 12/06
		Title: ptttype.Title_t{ //[心得] 然後呢？～
			0x5b, 0xa4, 0xdf, 0xb1, 0x6f, 0x5d, 0x20, 0xb5,
			0x4d, 0xab, 0xe1, 0xa9, 0x4f, 0xa1, 0x48, 0xa1,
			0xe3,
		},
		Filemode: ptttype.FILE_MULTI,
	}
	case_1_FileHeaders := []ptttype.FileHeaderRaw{
		*fileHeaderRaw1, // M.1607202239.A.30D
		*fileHeaderRaw2, // M.1607203395.A.00D
	}
	case_1_Filename := "./testcase/boards/W/WhoAmI/.DIR"
	defer os.RemoveAll(case_1_Filename)
	file, _ := os.OpenFile(case_1_Filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	defer file.Close()
	_ = types.BinaryWrite(file, binary.LittleEndian, case_1_FileHeaders)

	params0 := &DeleteArticlesParams{
		ArticleIDs: []bbs.ArticleID{bbs.ToArticleID(filename0)},
	}

	path0 := &DeleteArticlesPath{
		BBoardID: "10_WhoAmI",
	}

	type args struct {
		remoteAddr string
		uuserID    bbs.UUserID
		params     interface{}
		path       interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantResult interface{}
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			"",
			args{
				"127.0.0.1",
				"SYSOP",
				params0,
				path0,
			},
			DeleteArticlesResult{Indexes: []ptttype.SortIdx{1}},
			false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			wg.Done()
			gotResult, err := DeleteArticles(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DeleteArticles() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
		wg.Wait()
	}
}
