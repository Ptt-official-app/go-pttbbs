package api

import (
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

var testIP = "127.0.0.1"

func setupTest(name string) {
	types.SetIsTest()
	ptttype.SetIsTest()

	cache.SetIsTest()
	cmbbs.SetIsTest()

	err := types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")
	logrus.Infof("%v: after copy .PASSWDS: e: %v", name, err)

	err = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")
	logrus.Infof("%v: after copy .BRD: e: %v", name, err)

	err = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")
	logrus.Infof("%v: after copy boards: e: %v", name, err)

	err = types.CopyDirToDir("./testcase/home1", "./testcase/home")
	logrus.Infof("%v: after copy home: e: %v", name, err)

	time.Sleep(10 * time.Millisecond)

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()

	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest(name string) {
	defer time.Sleep(1 * time.Millisecond)

	defer func() {
		types.UnsetIsTest()
	}()

	defer ptttype.UnsetIsTest()

	defer cache.UnsetIsTest()

	defer cmbbs.UnsetIsTest()

	defer func() {
		os.Remove("./testcase/.PASSWDS")
		logrus.Infof("%v: after remove .PASSWDS", name)
	}()

	defer func() {
		os.Remove("./testcase/.BRD")
		logrus.Infof("%v: after remove .BRD", name)
	}()

	defer func() {
		os.RemoveAll("./testcase/boards")
		logrus.Infof("%v: after remove boards", name)
	}()

	defer func() {
		os.RemoveAll("./testcase/home")
		logrus.Infof("%v: after remove home", name)
	}()

	defer os.Remove("./testcase/.fresh")

	defer os.Remove("./testcase/.post")

	defer cache.CloseSHM()

	defer cmbbs.PasswdDestroy()

	defer freeTestVars()
}
