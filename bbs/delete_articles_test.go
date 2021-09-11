package bbs

import "testing"

func TestDeleteArticles(t *testing.T) {
	type args struct {
		uuserID    UUserID
		bboardID   BBoardID
		articleIDs []ArticleID
		ip         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteArticles(tt.args.uuserID, tt.args.bboardID, tt.args.articleIDs, tt.args.ip); (err != nil) != tt.wantErr {
				t.Errorf("DeleteArticles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
