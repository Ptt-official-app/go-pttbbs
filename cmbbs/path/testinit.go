package path

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	IsTest    = false
	TestMutex sync.Mutex
)

func SetIsTest() {
	logrus.Infof("path.SetIsTest: to TestMutex.Lock")
	TestMutex.Lock()
	logrus.Infof("path.SetIsTest: after TestMutex.Lock")
	IsTest = true
}

func UnsetIsTest() {
	IsTest = false
	logrus.Infof("path.UnsetIsTest: to TestMutex.Unlock")
	TestMutex.Unlock()
	logrus.Infof("path.UnsetIsTest: after TestMutex.Unlock")
}
