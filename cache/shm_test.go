package cache

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TestNewSHM(t *testing.T) {
	shmSetupTest()
	defer shmTeardownTest()

	type args struct {
		key          types.Key_t
		isUseHugeTlb bool
		isCreate     bool
	}
	tests := []struct {
		name        string
		args        args
		wantVersion int32
		wantSize    int32
		isClose     bool
		wantErr     bool
	}{
		// TODO: Add test caseShm.
		{
			args: args{
				key:          TestShmKey,
				isUseHugeTlb: false,
				isCreate:     true,
			},
			isClose:     false,
			wantVersion: SHM_VERSION,
			wantSize:    int32(SHM_RAW_SZ),
		},
		{
			args: args{
				key:          TestShmKey,
				isUseHugeTlb: false,
				isCreate:     false,
			},
			isClose:     true,
			wantVersion: SHM_VERSION,
			wantSize:    int32(SHM_RAW_SZ),
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			err := NewSHM(tt.args.key, tt.args.isUseHugeTlb, tt.args.isCreate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSHM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			Shm.Reset()

			Shm.Raw.Version = Shm.Shm.Version
			Shm.Raw.Size = Shm.Shm.Size

			if !reflect.DeepEqual(Shm.Raw.Version, tt.wantVersion) {
				t.Errorf("NewSHM() version: %v expected: %v", Shm.Raw.Version, tt.wantVersion)
			}

			if !reflect.DeepEqual(Shm.Raw.Size, tt.wantSize) {
				t.Errorf("NewSHM() size: %v expected :%v", Shm.Raw.Size, tt.wantSize)
			}

			if tt.isClose {
				CloseSHM()
			} else {
				Shm = nil
			}
		})
		wg.Wait()
	}
}
