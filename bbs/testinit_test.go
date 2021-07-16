package bbs

import (
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func setupTest() {
	types.SetIsTest()
	ptttype.SetIsTest()

	cache.SetIsTest()
	cmbbs.SetIsTest()

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

	_ = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")

	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")

	time.Sleep(1 * time.Millisecond)

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()
	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest() {
	defer time.Sleep(1 * time.Millisecond)

	defer types.UnsetIsTest()

	defer ptttype.UnsetIsTest()

	defer cache.UnsetIsTest()

	defer cmbbs.UnsetIsTest()

	defer os.Remove("./testcase/.PASSWDS")

	defer os.Remove("./testcase/.BRD")

	defer os.RemoveAll("./testcase/boards")

	defer os.RemoveAll("./testcase/home")

	defer os.Remove("./testcase/.fresh")

	defer os.Remove("./testcase/.post")

	defer cache.CloseSHM()

	defer cmbbs.PasswdDestroy()

	defer freeTestVars()
}
