package cmbbs

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
)

func setupTest() {
	SetIsTest()
	cache.SetIsTest()

	types.SetIsTest()

	ptttype.SetIsTest()

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
	_ = cache.CloseSHM()

	os.Remove("./testcase/.fresh")

	os.RemoveAll("./testcase/home")
	os.RemoveAll("./testcase/boards")
	os.Remove("./testcase/.BRD")
	os.Remove("./testcase/.PASSWDS")

	ptttype.UnsetIsTest()

	types.UnsetIsTest()

	cache.UnsetIsTest()
	UnsetIsTest()

}
