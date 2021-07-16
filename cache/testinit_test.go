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

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")
	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

	time.Sleep(1 * time.Millisecond)
}

func teardownTest() {
	defer time.Sleep(1 * time.Millisecond)

	defer shmTeardownTest()

	defer CloseSHM()

	defer os.Remove("./testcase/.PASSWDS")
	defer os.Remove("./testcase/.BRD")
}

func shmSetupTest() {
	types.SetIsTest("cache")

	ptttype.SetIsTest()

	SetIsTest()
}

func shmTeardownTest() {
	defer types.UnsetIsTest("cache")

	defer ptttype.UnsetIsTest()

	defer UnsetIsTest()
}
