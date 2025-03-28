package cache

import (
	"reflect"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

func TestLoadUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := Init(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		return
	}
	defer CloseSHM()

	wantHashHead := [1 << ptttype.HASH_BITS]ptttype.UIDInStore{}
	wantNextInHash := [ptttype.MAX_USERS]ptttype.UIDInStore{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 // SYSOP
	wantHashHead[56375] = 1 // CodingMan
	wantHashHead[36994] = 2 // pichu
	wantHashHead[15845] = 3 // Kahou
	wantHashHead[22901] = 4 // Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = ptttype.UIDInStore(idx + 1)
	}
	wantNextInHash[49] = -1

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
			var err error
			if err = LoadUHash(); (err != nil) != tt.wantErr {
				t.Errorf("loadUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			hashHead := &SHM.Shm.HashHead
			nextInHash := SHM.Shm.NextInHash

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}

			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)
			}
		})
		wg.Wait()
	}
}

func Test_fillUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := Init(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		return
	}
	defer CloseSHM()

	_ = LoadUHash()

	// move setupTest in for-loop
	wantHashHead := [1 << ptttype.HASH_BITS]ptttype.UIDInStore{}
	wantNextInHash := [ptttype.MAX_USERS]ptttype.UIDInStore{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 // SYSOP
	wantHashHead[56375] = 1 // CodingMan
	wantHashHead[36994] = 2 // pichu
	wantHashHead[15845] = 3 // Kahou
	wantHashHead[22901] = 4 // Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = ptttype.UIDInStore(idx + 1)
	}
	wantNextInHash[49] = -1

	type args struct {
		isOnfly bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "isOnfly",
			args: args{true},
		},
		{
			args: args{false},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			if err := fillUHash(tt.args.isOnfly); (err != nil) != tt.wantErr {
				t.Errorf("fillUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			hashHead := &SHM.Shm.HashHead
			nextInHash := SHM.Shm.NextInHash

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}
			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)
			}
		})
		wg.Wait()
	}
}

func TestInitFillUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	logrus.Infof("HASH_BITS: %v", ptttype.HASH_BITS)

	err := Init(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		return
	}
	defer CloseSHM()

	_ = LoadUHash()

	wantHashHead := [1 << ptttype.HASH_BITS]ptttype.UIDInStore{}
	wantNextInHash := [ptttype.MAX_USERS]ptttype.UIDInStore{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 // SYSOP
	wantHashHead[56375] = 1 // CodingMan
	wantHashHead[36994] = 2 // pichu
	wantHashHead[15845] = 3 // Kahou
	wantHashHead[22901] = 4 // Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = ptttype.UIDInStore(idx + 1)
	}
	wantNextInHash[49] = -1

	type args struct {
		isOnfly bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "isOnfly: true",
			args: args{true},
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()

			InitFillUHash(tt.args.isOnfly)

			hashHead := &SHM.Shm.HashHead
			nextInHash := SHM.Shm.NextInHash

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}

			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)
			}
		})
		wg.Wait()
	}
}
