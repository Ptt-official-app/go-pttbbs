package sem

import "sync"

const (
	testSemKey = 30000
)

var (
	testMutex sync.Mutex
)

func setupTest() {
	testMutex.Lock()

}

func teardownTest() {
	testMutex.Unlock()
}
