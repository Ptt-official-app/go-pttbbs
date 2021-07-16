package ptt

import (
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func setupTest(name string) {
	logrus.Infof("%v: to types.set-is-test", name)
	types.SetIsTest()

	ptttype.SetIsTest()

	cache.SetIsTest()

	cmbbs.SetIsTest()

	logrus.Infof("%v: to copy .PASSWDS", name)
	err := types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")
	if err != nil {
		logrus.Warn("unable to copy .PASSWDS")
	}

	logrus.Infof("%v: to copy .BRD", name)
	err = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")
	if err != nil {
		logrus.Warn("unable to copy .BRD")
	}

	logrus.Infof("%v: to copy boards", name)
	err = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")
	if err != nil {
		logrus.Warn("unable to copy boards")
	}

	logrus.Infof("%v: to copy home", name)
	err = types.CopyDirToDir("./testcase/home1", "./testcase/home")
	if err != nil {
		logrus.Warn("unable to copy home")
	}

	time.Sleep(1 * time.Millisecond)

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()

	_ = cmbbs.PasswdInit()

	initVars()
}

func teardownTest(name string) {
	defer time.Sleep(1 * time.Millisecond)

	defer func() {
		types.UnsetIsTest()
		logrus.Infof("%v: after types.unset-is-test", name)
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
	defer func() {
		// logrus.Infof("%v to remove .fresh")
		os.Remove("./testcase/.fresh")
	}()
	defer func() {
		// logrus.Infof("%v to remove .post")
		os.Remove("./testcase/.post")
	}()

	defer cache.CloseSHM()

	defer cmbbs.PasswdDestroy()

	defer freeTestVars()
}
