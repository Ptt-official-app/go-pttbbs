package initgin

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	jww "github.com/spf13/jwalterweatherman"
)

var ()

func setupTest() {

	jww.SetLogOutput(os.Stderr)
	//jww.SetLogThreshold(jww.LevelDebug)
	//jww.SetStdoutThreshold(jww.LevelDebug)
	log.SetLevel(log.DebugLevel)

	cache.SetIsTest()
	cmbbs.SetIsTest()

	types.SetIsTest()
	ptttype.SetIsTest()

	log.Infof("setupTest: to initAllConfig: sem_key: %v", ptttype.PASSWDSEM_KEY)

	_ = InitAllConfig("./testcase/test.ini")

	gin.SetMode(gin.TestMode)

	_ = types.CopyFileToFile("./testcase/.PASSWDS1", "./testcase/.PASSWDS")
	_ = types.CopyFileToFile("./testcase/.BRD1", "./testcase/.BRD")

	_ = types.CopyDirToDir("./testcase/boards1", "./testcase/boards")
	_ = types.CopyDirToDir("./testcase/home1", "./testcase/home")

	_ = cache.NewSHM(types.Key_t(cache.TestShmKey), ptttype.USE_HUGETLB, true)
	_ = cache.AttachSHM()

	cache.Shm.Reset()

	_ = cache.LoadUHash()

	cache.ReloadBCache()

	_ = cmbbs.PasswdInit()

	initTestVars()
}

func teardownTest() {
	_ = cmbbs.PasswdDestroy()

	_ = cache.CloseSHM()

	os.Remove("./testcase/.post")
	os.Remove("./testcase/.fresh")
	os.RemoveAll("./testcase/home")
	os.RemoveAll("./testcase/boards")
	os.Remove("./testcase/.BRD")
	os.Remove("./testcase/.PASSWDS")

	ptttype.UnsetIsTest()
	types.UnsetIsTest()

	cmbbs.UnsetIsTest()
	cache.UnsetIsTest()
}
