package cache

import (
	"sync"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

var (
	TestShmKey = types.Key_t(65534)
)

var (
	IsTest = false

	origBBSHome = ""
	origShmKey  = types.Key_t(0)

	TestMutex     sync.Mutex
	TestEachMutex sync.Mutex
)

func setupTest() {
	SetIsTest()

	origBBSHome = ptttype.SetBBSHOME("./testcase")

	initTestCases()
}

func teardownTest() {
	ptttype.SetBBSHOME(origBBSHome)

	UnsetIsTest()
	time.Sleep(1 * time.Millisecond)
}

func SetIsTest() {
	TestMutex.Lock()
	IsTest = true
	origShmKey = ptttype.SHM_KEY
	ptttype.SHM_KEY = TestShmKey
}

func UnsetIsTest() {
	ptttype.SHM_KEY = origShmKey
	IsTest = false
	TestMutex.Unlock()
}
