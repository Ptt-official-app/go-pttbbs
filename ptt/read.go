package ptt

import (
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func DeleteArticles(boardID *ptttype.BoardID_t, filename *ptttype.Filename_t) error {
	// rename filename
	// path.SetBFile(boardID, ptttype.FN_DIR)
	filePath, err := path.SetBFile(boardID, ptttype.FN_DIR)
	err = os.Rename(filePath+"/"+filename.String(), filename.DeletedName())
	if err != nil {
		return err
	}
	// rename.Dir content
	// cmsys.DeleteRecord()

	return nil
}
