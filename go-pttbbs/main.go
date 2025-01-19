package main

import (
	"github.com/Ptt-official-app/go-pttbbs/initgin"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := initMain()
	if err != nil {
		log.Errorf("unable to initMain: e: %v", err)
		return
	}
	router, err := initgin.InitGin()
	if err != nil {
		return
	}

	_ = router.Run(HTTP_HOST)
}
