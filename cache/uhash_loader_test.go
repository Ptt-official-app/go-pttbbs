package cache

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

func TestLoadUHash(t *testing.T) {
	setupTest()
	defer teardownTest()

	err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		return
	}
	defer CloseSHM()

	wantHashHead := [1 << ptttype.HASH_BITS]int32{}
	wantNextInHash := [ptttype.MAX_USERS]int32{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 //SYSOP
	wantHashHead[56375] = 1 //CodingMan
	wantHashHead[36994] = 2 //pichu
	wantHashHead[15845] = 3 //Kahou
	wantHashHead[22901] = 4 //Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = int32(idx + 1)
	}
	wantNextInHash[49] = -1

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if err = LoadUHash(); (err != nil) != tt.wantErr {
				t.Errorf("loadUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			hashHead := [1 << ptttype.HASH_BITS]int32{}
			nextInHash := [ptttype.MAX_USERS]int32{}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.HashHead),
				unsafe.Sizeof(Shm.Raw.HashHead),
				unsafe.Pointer(&hashHead),
			)

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NextInHash),
				unsafe.Sizeof(Shm.Raw.NextInHash),
				unsafe.Pointer(&nextInHash),
			)

			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)

			}

		})
	}
}

func Test_fillUHash(t *testing.T) {
	wantHashHead := [1 << ptttype.HASH_BITS]int32{}
	wantNextInHash := [ptttype.MAX_USERS]int32{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 //SYSOP
	wantHashHead[56375] = 1 //CodingMan
	wantHashHead[36994] = 2 //pichu
	wantHashHead[15845] = 3 //Kahou
	wantHashHead[22901] = 4 //Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = int32(idx + 1)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
			if err != nil {
				return
			}
			defer CloseSHM()

			_ = LoadUHash()

			if err := fillUHash(tt.args.isOnfly); (err != nil) != tt.wantErr {
				t.Errorf("fillUHash() error = %v, wantErr %v", err, tt.wantErr)
			}

			hashHead := [1 << ptttype.HASH_BITS]int32{}
			nextInHash := [ptttype.MAX_USERS]int32{}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.HashHead),
				unsafe.Sizeof(Shm.Raw.HashHead),
				unsafe.Pointer(&hashHead),
			)

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NextInHash),
				unsafe.Sizeof(Shm.Raw.NextInHash),
				unsafe.Pointer(&nextInHash),
			)

			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)

			}
		})
	}
}

func TestInitFillUHash(t *testing.T) {
	wantHashHead := [1 << ptttype.HASH_BITS]int32{}
	wantNextInHash := [ptttype.MAX_USERS]int32{}
	for idx := range wantHashHead {
		wantHashHead[idx] = -1
	}
	wantHashHead[29935] = 0 //SYSOP
	wantHashHead[56375] = 1 //CodingMan
	wantHashHead[36994] = 2 //pichu
	wantHashHead[15845] = 3 //Kahou
	wantHashHead[22901] = 4 //Kahou2
	wantHashHead[35] = 5    //""

	wantNextInHash[0] = -1
	wantNextInHash[1] = -1
	wantNextInHash[2] = -1
	wantNextInHash[3] = -1
	wantNextInHash[4] = -1
	for idx := 5; idx < 49; idx++ {
		wantNextInHash[idx] = int32(idx + 1)
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
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			setupTest()
			defer teardownTest()

			logrus.Infof("HASH_BITS: %v", ptttype.HASH_BITS)

			err := NewSHM(TestShmKey, ptttype.USE_HUGETLB, true)
			if err != nil {
				return
			}
			defer CloseSHM()

			_ = LoadUHash()

			InitFillUHash(tt.args.isOnfly)

			hashHead := [1 << ptttype.HASH_BITS]int32{}
			nextInHash := [ptttype.MAX_USERS]int32{}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.HashHead),
				unsafe.Sizeof(Shm.Raw.HashHead),
				unsafe.Pointer(&hashHead),
			)

			for idx, each := range hashHead {
				if each != wantHashHead[idx] {
					t.Errorf("loadUHash() (%v) hashHead: %v expected: %v", idx, each, wantHashHead[idx])
					break
				}
			}

			Shm.ReadAt(
				unsafe.Offsetof(Shm.Raw.NextInHash),
				unsafe.Sizeof(Shm.Raw.NextInHash),
				unsafe.Pointer(&nextInHash),
			)

			if !reflect.DeepEqual(nextInHash, wantNextInHash) {
				t.Errorf("loadUHash() nextInHash: %v expected: %v", nextInHash, wantNextInHash)

			}

		})
	}
}
