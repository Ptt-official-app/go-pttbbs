package cache

import (
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func setupTest() {
	shmSetupTest()

	initTestCases()

	err := NewSHM(types.Key_t(TestShmKey), ptttype.USE_HUGETLB, true)
	if err != nil {
		return
	}

	Shm.Reset()

}

func teardownTest() {
	CloseSHM()

	shmTeardownTest()
}

func shmSetupTest() {
	SetIsTest()

	origBBSHome = ptttype.SetBBSHOME("./testcase")
}

func shmTeardownTest() {
	ptttype.SetBBSHOME(origBBSHome)

	UnsetIsTest()
	time.Sleep(1 * time.Millisecond)
}
