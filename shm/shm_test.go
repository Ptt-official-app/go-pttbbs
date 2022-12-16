package shm

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

type testStruct struct {
	A int32
	B testStruct2
}

const TEST_STRUCT_SZ = unsafe.Sizeof(testStruct{})

type testStruct2 struct {
	C [10]uint8
}

func TestCreateShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	type args struct {
		key           types.Key_t
		size          types.Size_t
		is_usehugetlb bool
	}
	tests := []struct {
		name        string
		args        args
		wantShmid   int
		wantShmaddr unsafe.Pointer
		wantIsNew   bool
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:      "first-good-shm-id",
			args:      args{testShmKey, 100, false},
			wantShmid: 0,
			wantIsNew: true,
		},
		{
			name:      "not first-good-shm-id",
			args:      args{testShmKey, 100, false},
			wantShmid: 0,
			wantIsNew: false,
		},
	}

	firstGoodShmID := 0
	var firstGoodShmaddr unsafe.Pointer
	defer CloseShm(firstGoodShmID, firstGoodShmaddr)

	var wg sync.WaitGroup
	for idx, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotShmid, gotShmaddr, gotIsNew, err := CreateShm(tt.args.key, tt.args.size, tt.args.is_usehugetlb)
			log.Infof("(%v/%v): after CreateShm: gotShmid: %v gotShmaddr: %v gotIsNew: %v e: %v", idx, tt.name, gotShmid, gotShmaddr, gotIsNew, err)

			if (err != nil) != tt.wantErr {
				t.Errorf("(%v/%v): CreateShm() error = %v, wantErr %v", idx, tt.name, err, tt.wantErr)
				return
			}

			log.Infof("shm_test.CreateShm: to check firstGoodShmID: %v %v", firstGoodShmID, firstGoodShmaddr)
			if firstGoodShmID == 0 {
				firstGoodShmID = gotShmid
				firstGoodShmaddr = gotShmaddr
			}

			if tt.wantShmid != 0 && gotShmid != tt.wantShmid {
				t.Errorf("%v: CreateShm() gotShmid = %v, expected %v", tt.name, gotShmid, tt.wantShmid)
			}
			if gotIsNew != tt.wantIsNew {
				t.Errorf("%v: CreateShm() gotIsNew = %v, expected %v", tt.name, gotIsNew, tt.wantIsNew)
			}
		})
		wg.Wait()
	}
}

func TestCloseShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	gotShmid, gotShmaddr, _, _ := CreateShm(testShmKey, 100, false)
	log.Infof("TestCloseShm: after CreateShm: gotShmid: %v gotShmaddr: %v", gotShmid, gotShmaddr)

	type args struct {
		shmid   int
		shmaddr unsafe.Pointer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			args: args{gotShmid, gotShmaddr},
		},
		{
			args:    args{gotShmid, gotShmaddr},
			wantErr: true,
		},
	}
	var wg sync.WaitGroup
	for idx, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			err := CloseShm(tt.args.shmid, tt.args.shmaddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("(%v/%v) CloseShm: e: %v wantErr: %v", idx, tt.name, err, tt.wantErr)
			}
		})
		wg.Wait()
	}
}

func TestOpenShm(t *testing.T) {
	setupTest()
	defer teardownTest()

	shmid, shmaddr, _, _ := CreateShm(testShmKey, 100, false)
	defer CloseShm(shmid, shmaddr)

	type args struct {
		key           types.Key_t
		size          types.Size_t
		is_usehugetlb bool
	}
	tests := []struct {
		name        string
		args        args
		wantShmid   int
		wantShmaddr unsafe.Pointer
		wantErr     bool
	}{
		// TODO: Add test cases.
		{
			name:      "same size",
			args:      args{testShmKey, 100, false},
			wantShmid: shmid,
		},
		{
			name:      "same size 2",
			args:      args{testShmKey, 100, false},
			wantShmid: shmid,
		},
		{
			name:      "diff size",
			args:      args{testShmKey, 500, false},
			wantShmid: -1,
			wantErr:   true,
		},
		{
			name:      "diff key",
			args:      args{testShmKey + 1, 500, false},
			wantShmid: -1,
			wantErr:   true,
		},
	}
	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			gotShmid, _, err := OpenShm(tt.args.key, tt.args.size, tt.args.is_usehugetlb)
			if (err != nil) != tt.wantErr {
				t.Errorf("OpenShm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotShmid != tt.wantShmid {
				t.Errorf("OpenShm() gotShmid = %v, expected %v", gotShmid, tt.wantShmid)
			}
		})
		wg.Wait()
	}
}
