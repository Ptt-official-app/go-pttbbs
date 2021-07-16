package cmbbs

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func setupTest() {
	types.SetIsTest("cmbbs")

	ptttype.SetIsTest()

	cache.SetIsTest()

	SetIsTest()

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

	_ = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")

	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")

	err := cache.NewSHM(cache.TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("setupTest: unable to NewSHM: e: %v", err)
		return
	}

	cache.Shm.Reset()

	_ = cache.LoadUHash()
}

func teardownTest() {
	defer types.UnsetIsTest("cmbbs")

	defer ptttype.UnsetIsTest()

	defer cache.UnsetIsTest()

	defer UnsetIsTest()

	defer os.Remove("./testcase/.PASSWDS")

	defer os.Remove("./testcase/.BRD")

	defer os.RemoveAll("./testcase/boards")

	defer os.RemoveAll("./testcase/home")

	defer os.Remove("./testcase/.fresh")

	defer cache.CloseSHM()
}
