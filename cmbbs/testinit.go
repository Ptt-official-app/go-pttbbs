package cmbbs

import (
	"sync"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

var (
	origBBSHOME = ""
	IsTest      = false
	TestMutex   sync.Mutex

	TestPASSWDSEM_KEY = 32763

	origPASSWDSEM_KEY = 0
)

func setupTest() {
	SetIsTest()
	cache.SetIsTest()

	origBBSHOME = ptttype.SetBBSHOME("./testcase")

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

	cache.UnsetIsTest()
	UnsetIsTest()
}

func SetIsTest() {
	TestMutex.Lock()
	IsTest = true
	origPASSWDSEM_KEY = ptttype.PASSWDSEM_KEY
	ptttype.PASSWDSEM_KEY = TestPASSWDSEM_KEY

	log.Infof("After set sem: TestPASSWDSEM_KEY: %v ptttype.PASSWDSEM_KEY: %v", TestPASSWDSEM_KEY, ptttype.PASSWDSEM_KEY)
}

func UnsetIsTest() {
	IsTest = false
	TestMutex.Unlock()
}
