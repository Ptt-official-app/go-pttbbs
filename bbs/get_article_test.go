package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetArticle(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1)

	type args struct {
		uuserID    UUserID
		bboardID   BBoardID
		articleID  ArticleID
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
				uuserID:   UUserID("SYSOP"),
				bboardID:  BBoardID("10_WhoAmI"),
				articleID: articleID,
			},
			expectedContent: testContent1,
			expectedMtime:   1607969504,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, gotMtime, err := GetArticle(tt.args.uuserID, tt.args.bboardID, tt.args.articleID, tt.args.retrieveTS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.expectedContent) {
				t.Errorf("GetArticle() gotContent = %v, want %v", gotContent, tt.expectedContent)
			}
			if !reflect.DeepEqual(gotMtime, tt.expectedMtime) {
				t.Errorf("GetArticle() gotMtime = %v, want %v", gotMtime, tt.expectedMtime)
			}
		})
	}
}
