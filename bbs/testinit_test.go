package bbs

import (
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/boardd"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func setupTest() {
	types.SetIsTest("bbs")
	ptttype.SetIsTest()

	cache.SetIsTest()

	path.SetIsTest()

	cmbbs.SetIsTest()

	boardd.SetIsTest()

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

	_ = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")

	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")

	time.Sleep(1 * time.Millisecond)

	_ = cache.Init(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.SHM.Reset()

	cache.MAP.Reset()

	_ = cache.LoadUHash()
	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest() {
	defer time.Sleep(1 * time.Millisecond)

	defer types.UnsetIsTest("bbs")

	defer ptttype.UnsetIsTest()

	defer cache.UnsetIsTest()

	defer path.UnsetIsTest()

	defer cmbbs.UnsetIsTest()

	defer boardd.UnsetIsTest()

	defer os.Remove("./testcase/.PASSWDS")

	defer os.Remove("./testcase/.BRD")

	defer os.RemoveAll("./testcase/boards")

	defer os.RemoveAll("./testcase/home")

	defer os.Remove("./testcase/.fresh")

	defer os.Remove("./testcase/.post")

	defer cache.CloseSHM()

	defer cache.CloseMAP()

	defer cmbbs.PasswdDestroy()

	defer freeTestVars()
}
