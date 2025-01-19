package main

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	jww "github.com/spf13/jwalterweatherman"
)

func setupTest() {
	jww.SetLogOutput(os.Stderr)
	// jww.SetLogThreshold(jww.LevelDebug)
	// jww.SetStdoutThreshold(jww.LevelDebug)
	log.SetLevel(log.DebugLevel)

	types.SetIsTest("main")

	ptttype.SetIsTest()

	cache.SetIsTest()

	path.SetIsTest()

	cmbbs.SetIsTest()

	log.Infof("setupTest: to initAllConfig: sem_key: %v shm_key: %v", ptttype.PASSWDSEM_KEY, ptttype.SHM_KEY)

	_ = initAllConfig("./testcase/test.ini")

	log.Infof("setupTest: after initAllConfig: sem_key: %v shm_key: %v", ptttype.PASSWDSEM_KEY, ptttype.SHM_KEY)

	gin.SetMode(gin.TestMode)

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")

	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()

	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()
}

func teardownTest() {
	defer types.UnsetIsTest("main")

	defer ptttype.UnsetIsTest()

	defer cache.UnsetIsTest()

	defer path.UnsetIsTest()

	defer cmbbs.UnsetIsTest()

	defer os.Remove("./testcase/.PASSWDS")
	defer os.RemoveAll("./testcase/home")

	defer os.Remove("./testcase/.fresh")

	defer cache.CloseSHM()

	defer cmbbs.PasswdDestroy()
}
