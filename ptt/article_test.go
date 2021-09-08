package ptt

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func Test_readContent(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	SetupNewUser(testNewPostUser1)

	boardID0 := &ptttype.BoardID_t{}
	copy(boardID0[:], []byte("WhoAmI"))
	ip0 := &ptttype.IPv4_t{}
	copy(ip0[:], []byte("127.0.0.1"))

	class0 := []byte("test")
	title0 := []byte("this is a test")
	fullTitle0 := ptttype.Title_t{}
	copy(fullTitle0[:], []byte("[test] this is a test"))
	owner0 := ptttype.Owner_t{}
	copy(owner0[:], []byte("A1"))

	uid0, _ := cache.DoSearchUserRaw(&testNewPostUser1.UserID, nil)
	content0 := [][]byte{[]byte("test1"), []byte("test2")}

	articleSummary, _ := NewPost(testNewPostUser1, uid0, boardID0, 10, class0, title0, content0, ip0, nil)

	filename0, _ := path.SetBFile(boardID0, articleSummary.Filename.String())
	file0, _ := os.Open(filename0)
	defer file0.Close()
	postContent0, _ := io.ReadAll(file0)

	hash0 := cmsys.FNV1_64_INIT
	hash0 = cmsys.Fnv64Buf(postContent0, len(postContent0), hash0)

	type args struct {
		filename   string
		retrieveTS types.Time4
		isHash     bool
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		expectedHash    cmsys.Fnv64_t
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args:            args{filename: filename0},
			expectedContent: postContent0,
		},
		{
			args:            args{filename: filename0, isHash: true},
			expectedContent: postContent0,
			expectedHash:    hash0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, _, gotHash, err := readContent(tt.args.filename, tt.args.retrieveTS, tt.args.isHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("readContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("readContent() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
			if !reflect.DeepEqual(gotHash, tt.expectedHash) {
				t.Errorf("readContent() gotHash = %v, want %v", gotHash, tt.expectedHash)
			}
		})
	}
}
