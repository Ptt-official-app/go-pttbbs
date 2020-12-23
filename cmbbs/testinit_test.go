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

	origBBSHOME = ptttype.SetBBSHOME("./testcase")

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

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

	os.Remove("./testcase/.PASSWDS")

	ptttype.SetBBSHOME(origBBSHOME)

	cache.UnsetIsTest()
	UnsetIsTest()
}
