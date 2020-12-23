package ptttype

import (
	"os"

	jww "github.com/spf13/jwalterweatherman"
)

func setupTest() {
	jww.SetLogOutput(os.Stderr)
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
}

func teardownTest() {
}
