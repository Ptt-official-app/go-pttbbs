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

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()
	_ = cache.LoadUHash()
	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest() {
	freeTestVars()

	_ = cmbbs.PasswdDestroy()

	_ = cache.CloseSHM()

	os.Remove("./testcase/.fresh")
	os.Remove("./testcase/.PASSWDS")
	os.RemoveAll("./testcase/home")

	ptttype.SetBBSHOME(origBBSHOME)

	cmbbs.UnsetIsTest()
	cache.UnsetIsTest()
}
