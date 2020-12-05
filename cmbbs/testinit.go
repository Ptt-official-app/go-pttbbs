package cmbbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

var (
	origBBSHOME = ""
	IsTest      = false
)

func setupTest() {
	cache.TestMutex.Lock()
	origBBSHOME = ptttype.SetBBSHOME("./testcase")

	cache.IsTest = true
	err := cache.NewSHM(cache.TestShmKey, ptttype.USE_HUGETLB, true)
	if err != nil {
		log.Errorf("setupTest: unable to NewSHM: e: %v", err)
		return
	}
	_ = cache.LoadUHash()
}

func teardownTest() {
	_ = cache.CloseSHM()
	ptttype.SetBBSHOME(origBBSHOME)
	cache.IsTest = false
	cache.TestMutex.Unlock()
}
