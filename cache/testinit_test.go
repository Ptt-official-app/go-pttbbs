package cache

import (
	"os"
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

	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

}

func teardownTest() {
	os.Remove("./testcase/.BRD")

	CloseSHM()

	shmTeardownTest()
}

func shmSetupTest() {
	SetIsTest()

	ptttype.SetIsTest()

}

func shmTeardownTest() {
	ptttype.UnsetIsTest()

	UnsetIsTest()
	time.Sleep(1 * time.Millisecond)
}
