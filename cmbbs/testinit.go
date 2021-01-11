package cmbbs

import (
	"sync"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

var (
	IsTest    = false
	TestMutex sync.Mutex

	TestPASSWDSEM_KEY = 32763

	origPASSWDSEM_KEY = 0
)

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
