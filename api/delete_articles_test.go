package api

import (
	"reflect"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

func TestDeleteArticles(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := DeleteArticles(tt.args.remoteAddr, tt.args.uuserID, tt.args.params, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DeleteArticles() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
