package main

import (
	"flag"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	apiPrefix = "/v1"
)

var ()

func withPrefix(path string) string {
	return apiPrefix + path
}

func initGin() (*gin.Engine, error) {
	router := gin.Default()

	router.POST(
		withPrefix(api.INDEX_R),
		NewLoginRequiredApi(
			api.Index,
			&api.IndexParams{},
		).Json,
	)
	router.POST(
		withPrefix(api.LOGIN_R),
		NewApi(
			api.Login,
			&api.LoginParams{},
		).Json,
	)
	router.POST(
		withPrefix(api.REGISTER_R),
		NewApi(
			api.Register,
			&api.RegisterParams{},
		).Json,
	)
	router.GET(
		withPrefix(api.LOAD_GENERAL_BOARDS_R),
		NewLoginRequiredApi(
			api.LoadGeneralBoards,
			&api.LoadGeneralBoardsParams{},
		).Query,
	)
	router.GET(
		withPrefix(api.LOAD_GENERAL_ARTICLES_R),
		NewLoginRequiredPathApi(
			api.LoadGeneralArticles,
			&api.LoadGeneralArticlesParams{},
			&api.LoadGeneralArticlesPath{},
		).Query,
	)

	return router, nil
}

//initConfig
//
//Params
//	filename: ini filename
//
//Return
//	error: err
func initAllConfig(filename string) error {

	filenameList := strings.Split(filename, ".")
	if len(filenameList) == 1 {
		return ErrInvalidIni
	}

	filenamePrefix := strings.Join(filenameList[:len(filenameList)-1], ".")
	filenamePostfix := filenameList[len(filenameList)-1]
	viper.SetConfigName(filenamePrefix)
	viper.SetConfigType(filenamePostfix)
	viper.AddConfigPath("/etc/go-pttbbs")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	log.Debugf("viper keys: %v", viper.AllKeys())

	err = InitConfig()
	if err != nil {
		return err
	}
	err = api.InitConfig()
	if err != nil {
		return err
	}
	err = types.InitConfig()
	if err != nil {
		return err
	}
	err = ptttype.InitConfig()
	if err != nil {
		return err
	}
	err = ptt.InitConfig()
	if err != nil {
		return err
	}

	return InitConfig()
}

func initMain() error {
	jww.SetLogThreshold(jww.LevelDebug)
	jww.SetStdoutThreshold(jww.LevelDebug)
	log.SetLevel(log.InfoLevel)

	filename := ""
	flag.StringVar(&filename, "ini", "config.ini", "ini filename")
	flag.Parse()

	err := initAllConfig(filename)
	if err != nil {
		return err
	}

	//init shm
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

	//init sem
	err = cmbbs.PasswdInit()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := initMain()
	if err != nil {
		log.Errorf("unable to initMain: e: %v", err)
		return
	}
	router, err := initGin()
	if err != nil {
		return
	}

	_ = router.Run(HTTP_HOST)
}
