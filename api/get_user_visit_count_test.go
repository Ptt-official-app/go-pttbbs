package api

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/cache"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestGetUserVisitCount(t *testing.T) {
	setupTest(t.Name())
	defer teardownTest(t.Name())

	nUser := 5
	cache.Shm.WriteAt(
		unsafe.Offsetof(cache.Shm.Raw.UTMPNumber),
		types.INT32_SZ,
		unsafe.Pointer(&nUser),
	)

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserVisitCount(testIP, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserVisitCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserVisitCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
