package bbs

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestLoadGeneralArticlesAllGuest(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		bboardID    BBoardID
		startIdxStr string
		nArticles   int
		isDesc      bool
	}
	tests := []struct {
		name               string
		args               args
		wantSummaries      []*ArticleSummary
		wantNextIdxStr     string
		wantNextCreateTime types.Time4
		wantIsNewest       bool
		wantErr            bool
	}{
		// TODO: Add test cases.
		{
			name: "most basic setup: find 1 article with desc, starting from nil",
			args: args{
				bboardID:  "WhoAmI",
				nArticles: 1,
				isDesc:    true,
			},
			wantSummaries:      []*ArticleSummary{testArticleSummary1AllGuest},
			wantNextIdxStr:     "1607202239@1Vo_M_CD",
			wantNextCreateTime: 1607202239,
			wantIsNewest:       true,
		},
		{
			name: "most basic setup: find 2 articles with desc, starting from nil",
			args: args{
				bboardID:  "WhoAmI",
				nArticles: 2,
				isDesc:    true,
			},
			wantSummaries: []*ArticleSummary{testArticleSummary1AllGuest, testArticleSummary0AllGuest},
			wantIsNewest:  true,
		},
		{
			name: "most basic setup: find 1 article with asc (not desc), starting from nil",
			args: args{
				bboardID:  "WhoAmI",
				nArticles: 1,
				isDesc:    false,
			},
			wantSummaries:      []*ArticleSummary{testArticleSummary0AllGuest},
			wantNextIdxStr:     "1607203395@1Vo_f30D",
			wantNextCreateTime: 1607203395,
		},
		{
			name: "most basic setup: find 2 articles with asc (not desc), starting from nil",
			args: args{
				bboardID:  "WhoAmI",
				nArticles: 2,
				isDesc:    false,
			},
			wantSummaries: []*ArticleSummary{testArticleSummary0AllGuest, testArticleSummary1AllGuest},
			wantIsNewest:  true,
		},
		{
			name: "most basic setup: find 1 article with desc, starting from the last article",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607203395@1Vo_f30D",
				nArticles:   1,
				isDesc:      true,
			},
			wantSummaries:      []*ArticleSummary{testArticleSummary1AllGuest},
			wantNextIdxStr:     "1607202239@1Vo_M_CD",
			wantNextCreateTime: 1607202239,
			wantIsNewest:       true,
		},
		{
			name: "find 2 articles with asc (not isDesc), starting from the last article, should return 1 article",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607203395@1Vo_f30D",
				nArticles:   2,
				isDesc:      false,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary1AllGuest},
			wantNextIdxStr: "",
			wantIsNewest:   true,
		},
		{
			name: "find 1 article with asc (not isDesc), starting from the beginning article",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CD",
				nArticles:   1,
				isDesc:      false,
			},
			wantSummaries:      []*ArticleSummary{testArticleSummary0AllGuest},
			wantNextIdxStr:     "1607203395@1Vo_f30D",
			wantNextCreateTime: 1607203395,
			wantIsNewest:       false,
		},
		{
			name: "find 2 articles with asc (not desc), starting from the beginning article, should return 2 articles",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CD",
				nArticles:   2,
				isDesc:      false,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary0AllGuest, testArticleSummary1AllGuest},
			wantNextIdxStr: "",
			wantIsNewest:   true,
		},
		{
			name: "find 2 articles with asc (not desc), starting from same create-time of the beginning article, but with diff filename, should return 1 articles",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CE",
				nArticles:   2,
				isDesc:      false,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary0AllGuest, testArticleSummary1AllGuest},
			wantNextIdxStr: "",
			wantIsNewest:   true,
		},
		{
			name: "find 2 articles with desc, starting from same create-time of the beginning article, but with diff filename, should return err.",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202239@1Vo_M_CE",
				nArticles:   2,
				isDesc:      true,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary0AllGuest},
			wantNextIdxStr: "",
		},
		{
			name: "find 2 articles with asc (not desc), starting after the beginning article, should return 1 article (last article)",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202240@1Vo_N0CD",
				nArticles:   2,
				isDesc:      false,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary1AllGuest},
			wantNextIdxStr: "",
			wantIsNewest:   true,
		},
		{
			name: "find 2 articles with desc, starting after the beginning article, should return 1 article (beginning article)",
			args: args{
				bboardID:    "WhoAmI",
				startIdxStr: "1607202240@1Vo_N0CD",
				nArticles:   2,
				isDesc:      true,
			},
			wantSummaries:  []*ArticleSummary{testArticleSummary0AllGuest},
			wantNextIdxStr: "",
			wantIsNewest:   false,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, gotNextIdxStr, gotNextCreateTime, gotIsNewest, _, err := LoadGeneralArticlesAllGuest(tt.args.bboardID, tt.args.startIdxStr, tt.args.nArticles, tt.args.isDesc)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadGeneralArticlesAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries, tt.wantSummaries) {
				t.Errorf("LoadGeneralArticlesAllGuest() gotSummaries = %v, want %v", gotSummaries, tt.wantSummaries)
			}
			if gotNextIdxStr != tt.wantNextIdxStr {
				t.Errorf("LoadGeneralArticlesAllGuest() gotNextIdxStr = %v, want %v", gotNextIdxStr, tt.wantNextIdxStr)
			}
			if !reflect.DeepEqual(gotNextCreateTime, tt.wantNextCreateTime) {
				t.Errorf("LoadGeneralArticlesAllGuest() gotNextCreateTime = %v, want %v", gotNextCreateTime, tt.wantNextCreateTime)
			}
			if gotIsNewest != tt.wantIsNewest {
				t.Errorf("LoadGeneralArticlesAllGuest() gotIsNewest = %v, want %v", gotIsNewest, tt.wantIsNewest)
			}
		})
	}
	wg.Wait()
}
