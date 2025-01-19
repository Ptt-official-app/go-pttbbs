package bbs

import (
	"reflect"
	"sync"
	"testing"
)

func TestLoadBottomArticlesAllGuest(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		bboardID BBoardID
	}
	tests := []struct {
		name          string
		args          args
		wantSummaries []*ArticleSummary
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			args: args{
				bboardID: "WhoAmI",
			},
			wantSummaries: []*ArticleSummary{testBottomSummary1AllGuest},
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotSummaries, err := LoadBottomArticlesAllGuest(tt.args.bboardID)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadBottomArticlesAllGuest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSummaries[0], tt.wantSummaries[0]) {
				t.Errorf("LoadBottomArticlesAllGuest() = %v, want %v", gotSummaries, tt.wantSummaries)
			}
		})
	}
	wg.Wait()
}
