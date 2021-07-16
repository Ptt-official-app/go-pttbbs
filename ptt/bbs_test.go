package ptt

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestReadPost(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.ReloadBCache()

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("WhoAmI"))

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))

	filename := "testcase/boards/W/WhoAmI/M.1607202239.A.30D"
	mtime := time.Unix(1607209066, 0)
	os.Chtimes(filename, mtime, mtime)

	filename2 := &ptttype.Filename_t{}
	copy(filename2[:], []byte("M.1607202239.A.31D"))
	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardID    *ptttype.BoardID_t
		bid        ptttype.Bid
		filename   *ptttype.Filename_t
		retrieveTS types.Time4
	}
	tests := []struct {
		name            string
		args            args
		expectedContent []byte
		expectedMtime   types.Time4
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:     testUserecRaw1,
				uid:      1,
				boardID:  boardID1,
				bid:      10,
				filename: filename1,
			},
			expectedContent: testContent1,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename1,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   1607209066,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID1,
				bid:        10,
				filename:   filename2,
				retrieveTS: 1607209066,
			},
			expectedContent: nil,
			expectedMtime:   0,
			wantErr:         true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotContent, gotMtime, err := ReadPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.filename, tt.args.retrieveTS)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("ReadPost() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.expectedMtime) {
				t.Errorf("ReadPost() gotMtime = %v, want %v", gotMtime, tt.expectedMtime)
			}
		})
	}
	wg.Wait()
}

func TestNewPost(t *testing.T) {
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
	expectedSummary0 := &ptttype.ArticleSummaryRaw{
		Aid:     3,
		BoardID: boardID0,
		Class:   class0,
		FileHeaderRaw: &ptttype.FileHeaderRaw{
			Title: fullTitle0,
			Owner: owner0,
		},
	}

	expectedSummary1 := &ptttype.ArticleSummaryRaw{
		Aid:     4,
		BoardID: boardID0,
		Class:   class0,
		FileHeaderRaw: &ptttype.FileHeaderRaw{
			Title: fullTitle0,
			Owner: owner0,
		},
	}

	nuser := int32(40)
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.BCache)+ptttype.BOARD_HEADER_RAW_SZ*uintptr(9)+ptttype.BOARD_HEADER_NUSER_OFFSET,
		types.INT32_SZ,
		unsafe.Pointer(&nuser),
	)

	uid0, _ := cache.DoSearchUserRaw(&testNewPostUser1.UserID, nil)

	content0 := [][]byte{[]byte("test1"), []byte("test2")}
	expected0 := []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 'A', '1', ' ', // 作者: A1
		0x28, 0xaf, 0xab, //(神
		0x29, 0x20, 0xac, 0xdd, 0xaa, 0x4f, //) 看板
		0x3a, 0x20, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49, 0x0a, //: WhoAmI
		0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0x74, 0x65, 0x73, // 標題: [tes
		0x74, 0x5d, 0x20, 0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, // t] this is
		0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74, 0x0a, // a test
		0xae, 0xc9, 0xb6, 0xa1, 0x3a, 0x20, 0x00, 0x00, 0x00, 0x20, // 時間: 000
		0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x20, 0x00, 0x00, 0x3a, // 000 00 00:
		0x00, 0x00, 0x3a, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x0a, // 00:00 0000
		0x74, 0x65, 0x73, 0x74, 0x31, 0x0a, // test1
		0x74, 0x65, 0x73, 0x74, 0x32, 0x0a, // test2
		0x0a,
		0x2d, 0x2d, 0x0a, //--
		0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, //※ 發信站:
		0x20, 0xb7, 0x73, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd, 0xf0, 0x28, // 新批踢踢(
		0x70, 0x74, 0x74, 0x32, 0x2e, 0x63, 0x63, 0x29, 0x2c, 0x20, // ptt2.cc),
		0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x32, 0x37, 0x2e, // 來自: 127.
		0x30, 0x2e, 0x30, 0x2e, 0x31, 0x0a, // 0.0.1
		0xa1, 0xb0, 0x20, 0xa4, 0xe5, 0xb3, 0xb9, 0xba, 0xf4, //※ 文章網
		0xa7, 0x7d, 0x3a, 0x20, 0x68, 0x74, 0x74, 0x70, 0x3a, 0x2f, // 址: http:/
		0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, ///localhost
		0x2f, 0x62, 0x62, 0x73, 0x2f, 0x57, 0x68, 0x6f, 0x41, 0x6d, ///bbs/WhoAm
		0x49, 0x2f, 0x4d, 0x2e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // I/M.000000
		0x00, 0x00, 0x00, 0x00, 0x2e, 0x41, 0x2e, 0x00, 0x00, 0x00, // 0000.A.000
		0x2e, 0x68, 0x74, 0x6d, 0x6c, 0x0a, //.html
	}

	removeIdxes := []int{
		61, 62, 63, 65, 66, 67, 69, 70, 72, 73, 75, 76, 78, 79, 81, 82, 83, 84, // 時間
		191, 192, 193, 194, 195, 196, 197, 198, 199, 200, 204, 205, 206, // 文章網址
	}

	type args struct {
		user     *ptttype.UserecRaw
		uid      ptttype.Uid
		boardID  *ptttype.BoardID_t
		bid      ptttype.Bid
		posttype []byte
		title    []byte
		content  [][]byte
		ip       *ptttype.IPv4_t
		from     []byte
	}
	tests := []struct {
		name            string
		args            args
		expectedSummary *ptttype.ArticleSummaryRaw
		expected        []byte
		wantErr         bool
	}{
		// TODO: Add test cases.
		{
			name: "post0",
			args: args{
				user:     testNewPostUser1,
				uid:      uid0,
				boardID:  boardID0,
				bid:      10,
				posttype: class0,
				title:    title0,
				content:  content0,
				ip:       ip0,
			},
			expectedSummary: expectedSummary0,
			expected:        expected0,
		},
		{
			name: "post1",
			args: args{
				user:     testNewPostUser1,
				uid:      uid0,
				boardID:  boardID0,
				bid:      10,
				posttype: class0,
				title:    title0,
				content:  content0,
				ip:       ip0,
			},
			expectedSummary: expectedSummary1,
			expected:        expected0,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummary, err := NewPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.posttype, tt.args.title, tt.args.content, tt.args.ip, tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			content, mtime, err := ReadPost(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, &gotSummary.FileHeaderRaw.Filename, 0)
			if err != nil {
				t.Errorf("NewPost() unable to ReadPost: e: %v", err)
				return
			}

			if mtime != gotSummary.Modified {
				t.Errorf("NewPost() mtime: %v expected: %v", mtime, gotSummary.Modified)
			}
			gotSummary.Filename = ptttype.Filename_t{}
			gotSummary.Modified = 0
			gotSummary.Date = ptttype.Date_t{}
			testutil.TDeepEqual(t, "summary", gotSummary, tt.expectedSummary)

			for _, idx := range removeIdxes {
				if idx >= len(content) {
					break
				}
				content[idx] = 0x00
			}

			testutil.TDeepEqual(t, "content", content, tt.expected)
		})
		wg.Wait()
	}
}
