package bbs

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

var (
	origBBSHOME string
)

func setupTest() {
	cache.SetIsTest()
	cmbbs.SetIsTest()

	origBBSHOME = ptttype.SetBBSHOME("./testcase")

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")
	_ = types.CopyDirToDir("./testcase/board1", "./testcase/boards")

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()
	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest() {
	freeTestVars()

	os.Remove("./testcase/.fresh")

	_ = cmbbs.PasswdDestroy()

	_ = cache.CloseSHM()

	os.RemoveAll("./testcase/boards")
	os.RemoveAll("./testcase/home")
	os.Remove("./testcase/.PASSWDS")

	ptttype.SetBBSHOME(origBBSHOME)

	cmbbs.UnsetIsTest()
	cache.UnsetIsTest()
}
