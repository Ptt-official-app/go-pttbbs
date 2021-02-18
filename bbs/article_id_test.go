package bbs

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestToArticleID(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607937174.A.081"))
	articleID := ArticleID("1VrooM21SYSOP")

	type args struct {
		filename *ptttype.Filename_t
		ownerID  UUserID
	}
	tests := []struct {
		name     string
		args     args
		expected ArticleID
	}{
		// TODO: Add test cases.
		{
			args:     args{filename: filename1, ownerID: UUserID("SYSOP")},
			expected: articleID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToArticleID(tt.args.filename, tt.args.ownerID); got != tt.expected {
				t.Errorf("ToArticleID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestArticleID_ToRaw(t *testing.T) {
	setupTest()
	defer teardownTest()

	filename1 := &ptttype.Filename_t{}
	copy(filename1[:], []byte("M.1607202239.A.30D"))
	articleID := ToArticleID(filename1, "SYSOP")

	tests := []struct {
		name             string
		a                ArticleID
		expectedFilename *ptttype.Filename_t
		expectedOwnerID  UUserID
		wantErr          bool
	}{
		// TODO: Add test cases.
		{
			a:                articleID,
			expectedFilename: filename1,
			expectedOwnerID:  "SYSOP",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilename, gotOwnerID := tt.a.ToRaw()
			if !reflect.DeepEqual(gotFilename, tt.expectedFilename) {
				t.Errorf("ArticleID.ToRaw() gotFilename = %v, want %v", gotFilename, tt.expectedFilename)
			}
			if gotOwnerID != tt.expectedOwnerID {
				t.Errorf("ArticleID.ToRaw() gotOwnerID = %v, want %v", gotOwnerID, tt.expectedOwnerID)
			}
		})
	}
}
