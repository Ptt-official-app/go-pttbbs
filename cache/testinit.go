package cache

import (
	"sync"

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
