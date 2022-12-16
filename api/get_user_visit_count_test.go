package api

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/cache"
)

func TestGetUserVisitCount(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	cache.Shm.Shm.UTMPNumber = 5

	expected := &GetUserVisitCountResult{5}
	tests := []struct {
		name    string
		want    interface{}
		wantErr bool
	}{
		{
			"test get user visit count",
			expected,
			false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			got, err := GetUserVisitCount(testIP, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserVisitCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserVisitCount() got = %v, want %v", got, tt.want)
			}
		})
		wg.Wait()
	}
}
