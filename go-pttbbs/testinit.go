package main

import (
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	jww "github.com/spf13/jwalterweatherman"
)

func setupTest() {
	jww.SetLogOutput(os.Stderr)
	//jww.SetLogThreshold(jww.LevelDebug)
	//jww.SetStdoutThreshold(jww.LevelDebug)
	log.SetLevel(log.DebugLevel)

	_ = initAllConfig("./testcase/test.ini")

	gin.SetMode(gin.TestMode)

}

func teardownTest() {
}
