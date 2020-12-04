package cache

import (
	"sync"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

var (
	TestShmKey = types.Key_t(1000000)
)

var (
	IsTest = false

	origBBSHome = ""
	origShmKey  = types.Key_t(0)

	TestMutex     sync.Mutex
	TestEachMutex sync.Mutex
)

func setupTest() {
	TestMutex.Lock()
	IsTest = true
	origShmKey = ptttype.SHM_KEY
	ptttype.SHM_KEY = TestShmKey
	origBBSHome = ptttype.SetBBSHOME("./testcase")

	initTestCases()
}

func teardownTest() {
	ptttype.SetBBSHOME(origBBSHome)
	ptttype.SHM_KEY = origShmKey
	IsTest = false
	TestMutex.Unlock()
	time.Sleep(1 * time.Millisecond)
}
