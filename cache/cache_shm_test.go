package cache

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func TestAttachSHM(t *testing.T) {
	setupTest()
	defer teardownTest()

	log.Infof("TestAttachSHM: to NewSHM: shm_key: %v USE_HUGETLB: %v", ptttype.SHM_KEY, ptttype.USE_HUGETLB)

	log.Infof("TestAttachSHM: after NewSHM")

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := AttachSHM(); (err != nil) != tt.wantErr {
				t.Errorf("AttachSHM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	wg.Wait()
}

func TestAttachCheckSHM(t *testing.T) {
	setupTest()
	defer teardownTest()

	loaded := int32(1)
	Shm.WriteAt(
		unsafe.Offsetof(Shm.Raw.Loaded),
		types.INT32_SZ,
		unsafe.Pointer(&loaded),
	)

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			wantErr: false,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			if err := AttachCheckSHM(); (err != nil) != tt.wantErr {
				t.Errorf("AttachCheckSHM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}
