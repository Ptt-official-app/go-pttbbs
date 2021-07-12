package ptt

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticles(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	bid := ptttype.Bid(10)

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))
	bid1 := ptttype.Bid(1)

	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardIDRaw *ptttype.BoardID_t
		bid        ptttype.Bid
		startAid   ptttype.SortIdx
		nArticles  int
		isDesc     bool
	}
	tests := []struct {
		name                string
		args                args
		expectedSummaries   []*ptttype.ArticleSummaryRaw
		expectedIsNewest    bool
		expectedNextSummary *ptttype.ArticleSummaryRaw
		expectedStartNumIdx ptttype.SortIdx
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  1,
				isDesc:     true,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedNextSummary: testArticleSummary0,
			expectedStartNumIdx: 2,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  100,
				isDesc:     true,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1, testArticleSummary0},
			expectedStartNumIdx: 2,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  2,
				isDesc:     true,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1, testArticleSummary0},
			expectedStartNumIdx: 2,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  1,
				isDesc:     true,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedIsNewest:    true,
			expectedStartNumIdx: 2,
			expectedNextSummary: testArticleSummary0,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   1,
				nArticles:  1,
				isDesc:     true,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary0},
			expectedStartNumIdx: 1,
			expectedIsNewest:    false,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID1,
				bid:        bid1,
				startAid:   2,
				nArticles:  1,
				isDesc:     true,
			},
			expectedSummaries: nil,
			expectedIsNewest:  true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  1,
				isDesc:     false,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedNextSummary: nil,
			expectedStartNumIdx: 2,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   1,
				nArticles:  100,
				isDesc:     false,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary0, testArticleSummary1},
			expectedStartNumIdx: 1,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   1,
				nArticles:  2,
				isDesc:     false,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary0, testArticleSummary1},
			expectedStartNumIdx: 1,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   1,
				nArticles:  1,
				isDesc:     false,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary0},
			expectedIsNewest:    false,
			expectedStartNumIdx: 1,
			expectedNextSummary: testArticleSummary1,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
				startAid:   2,
				nArticles:  1,
				isDesc:     false,
			},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedStartNumIdx: 2,
			expectedIsNewest:    true,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID1,
				bid:        bid1,
				startAid:   1,
				nArticles:  1,
				isDesc:     false,
			},
			expectedSummaries: nil,
			expectedIsNewest:  true,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotIsNewest, gotNextSummary, gotStartNumIdx, err := LoadGeneralArticles(tt.args.user, tt.args.uid, tt.args.boardIDRaw, tt.args.bid, tt.args.startAid, tt.args.nArticles, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, each := range gotSummaries {
				each.BoardID = nil
			}

			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)

			if gotNextSummary != nil {
				gotNextSummary.BoardID = nil
			}
			testutil.TDeepEqual(t, "nextSummary", gotNextSummary, tt.expectedNextSummary)

			if !reflect.DeepEqual(gotIsNewest, tt.expectedIsNewest) {
				t.Errorf("LoadGeneralArticles() isNewest = %v, want %v", gotIsNewest, tt.expectedIsNewest)
			}

			if gotStartNumIdx != tt.expectedStartNumIdx {
				t.Errorf("LoadGeneralArticles() startNumIdx = %v, want %v", gotStartNumIdx, tt.expectedStartNumIdx)

			}
		})
	}
	wg.Wait()
}

func TestFindArticleStartAid(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	bid := ptttype.Bid(10)

	/*
		boardID1 := &ptttype.BoardID_t{}
		copy(boardID1[:], []byte("SYSOP"))
		bid1 := ptttype.Bid(1)
	*/

	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardID    *ptttype.BoardID_t
		bid        ptttype.Bid
		createTime types.Time4
		filename   *ptttype.Filename_t
		isDesc     bool
	}
	tests := []struct {
		name             string
		args             args
		expectedStartAid ptttype.SortIdx
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607202239,
				filename:   &testArticleSummary0.FileHeaderRaw.Filename,
				isDesc:     false,
			},
			expectedStartAid: 1,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607203395,
				filename:   &testArticleSummary1.FileHeaderRaw.Filename,
				isDesc:     false,
			},
			expectedStartAid: 2,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607202239,
				filename:   &testArticleSummary0.FileHeaderRaw.Filename,
				isDesc:     true,
			},
			expectedStartAid: 1,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607203395,
				filename:   &testArticleSummary1.FileHeaderRaw.Filename,
				isDesc:     true,
			},
			expectedStartAid: 2,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607202239,
				filename:   nil,
				isDesc:     true,
			},
			expectedStartAid: -1,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607203395,
				filename:   nil,
				isDesc:     true,
			},
			expectedStartAid: 1,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607202239,
				filename:   nil,
				isDesc:     false,
			},
			expectedStartAid: 2,
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardID:    boardID,
				bid:        bid,
				createTime: 1607203395,
				filename:   nil,
				isDesc:     false,
			},
			expectedStartAid: -1,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotStartAid, err := FindArticleStartIdx(tt.args.user, tt.args.uid, tt.args.boardID, tt.args.bid, tt.args.createTime, tt.args.filename, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindArticleStartAid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStartAid, tt.expectedStartAid) {
				t.Errorf("FindArticleStartAid() = %v, want %v", gotStartAid, tt.expectedStartAid)
			}
		})
	}
	wg.Wait()
}

func TestLoadGeneralArticlesSameCreateTime(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	bid := ptttype.Bid(10)

	type args struct {
		boardIDRaw *ptttype.BoardID_t
		bid        ptttype.Bid
		startAid   ptttype.SortIdx
		endAid     ptttype.SortIdx
		createTime types.Time4
	}
	tests := []struct {
		name                string
		args                args
		expectedSummaries   []*ptttype.ArticleSummaryRaw
		expectedStartNumIdx ptttype.SortIdx
		expectedEndNumIdx   ptttype.SortIdx
		wantErr             bool
	}{
		// TODO: Add test cases.
		{
			args:                args{boardIDRaw: boardID, bid: bid, startAid: 1, endAid: 0, createTime: 1607202239},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary0},
			expectedStartNumIdx: 1,
			expectedEndNumIdx:   1,
		},
		{
			args:                args{boardIDRaw: boardID, bid: bid, startAid: 1, endAid: 0, createTime: 1607203395},
			expectedSummaries:   []*ptttype.ArticleSummaryRaw{testArticleSummary1},
			expectedStartNumIdx: 2,
			expectedEndNumIdx:   2,
		},
		{
			args:              args{boardIDRaw: boardID, bid: bid, startAid: 1, endAid: 0, createTime: 1607203394},
			expectedSummaries: nil,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotStartNumIdx, gotEndNumIdx, err := LoadGeneralArticlesSameCreateTime(tt.args.boardIDRaw, tt.args.bid, tt.args.startAid, tt.args.endAid, tt.args.createTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticlesSameCreateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, each := range gotSummaries {
				each.BoardID = nil
			}
			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)

			if gotStartNumIdx != tt.expectedStartNumIdx {
				t.Errorf("LoadGeneralArticlesSameCreateTime() startNumIdx = %v, want %v", gotStartNumIdx, tt.expectedStartNumIdx)

			}

			if gotEndNumIdx != tt.expectedEndNumIdx {
				t.Errorf("LoadGeneralArticlesSameCreateTime() endNumIdx = %v, want %v", gotEndNumIdx, tt.expectedEndNumIdx)

			}
		})
	}
	wg.Wait()
}

func TestLoadBottomArticles(t *testing.T) {
	setupTest()
	defer teardownTest()

	cache.ReloadBCache()

	boardID := &ptttype.BoardID_t{}
	copy(boardID[:], []byte("WhoAmI"))
	bid := ptttype.Bid(10)

	boardID1 := &ptttype.BoardID_t{}
	copy(boardID1[:], []byte("SYSOP"))
	bid1 := ptttype.Bid(1)

	type args struct {
		user       *ptttype.UserecRaw
		uid        ptttype.Uid
		boardIDRaw *ptttype.BoardID_t
		bid        ptttype.Bid
	}
	tests := []struct {
		name              string
		args              args
		expectedSummaries []*ptttype.ArticleSummaryRaw
		wantErr           bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID,
				bid:        bid,
			},
			expectedSummaries: []*ptttype.ArticleSummaryRaw{testBottomSummary1},
		},
		{
			args: args{
				user:       testUserecRaw1,
				uid:        1,
				boardIDRaw: boardID1,
				bid:        bid1,
			},
			expectedSummaries: nil,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadBottomArticles(tt.args.user, tt.args.uid, tt.args.boardIDRaw, tt.args.bid)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBottomArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, each := range gotSummaries {
				each.BoardID = nil
				copy(each.Multi[:], []byte{0, 0, 0, 0})
			}

			testutil.TDeepEqual(t, "summaries", gotSummaries, tt.expectedSummaries)
		})
	}
	wg.Wait()
}
