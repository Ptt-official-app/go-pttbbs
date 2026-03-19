package initgin

import (
	"github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/boardd"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/configutil"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// initConfig
//
// Params
//
//	filename: ini filename
//
// Return
//
//	error: err
func InitAllConfig(filename string) (err error) {
	err = configutil.InitViper(filename)
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

	err = boardd.InitConfig()
	if err != nil {
		return err
	}

	err = path.InitConfig()
	if err != nil {
		return err
	}

	return nil
}
