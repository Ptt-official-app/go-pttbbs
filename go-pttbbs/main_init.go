package main

import (
	"flag"

	"github.com/Ptt-official-app/go-pttbbs/boardd"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/initgin"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
)

func initAllConfig(filename string) (err error) {
	err = initgin.InitAllConfig(filename)
	if err != nil {
		return err
	}

	return InitConfig()
}

func initMain() error {
	jww.SetLogThreshold(jww.LevelInfo)
	jww.SetStdoutThreshold(jww.LevelInfo)
	log.SetLevel(log.InfoLevel)

	filename := ""
	flag.StringVar(&filename, "ini", "config.ini", "ini filename")
	flag.Parse()

	err := initAllConfig(filename)
	if err != nil {
		return err
	}

	// init shm
	err = cache.NewSHM(types.Key_t(ptttype.SHM_KEY), ptttype.USE_HUGETLB, ptttype.IS_NEW_SHM)
	if err != nil {
		log.Errorf("unable to init SHM: e: %v", err)
		return err
	}

	if ptttype.IS_NEW_SHM {
		err = cache.LoadUHash()
		if err != nil {
			log.Errorf("unable to load UHash: e: %v", err)
			return err
		}

		cache.ReloadBCache()
	}
	err = cache.AttachCheckSHM()
	if err != nil {
		log.Errorf("unable to attach-check-shm: e: %v", err)
		return err
	}

	// init sem
	err = cmbbs.PasswdInit()
	if err != nil {
		return err
	}

	// init grpc
	err = boardd.Init(false)
	if err != nil {
		return err
	}

	return nil
}
