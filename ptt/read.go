package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func DeleteArticles(boardID *ptttype.BoardID_t, filename *ptttype.Filename_t, index ptttype.SortIdx) error {
	// path.SetBFile(boardID, ptttype.FN_DIR)
	filePath, err := path.SetBFile(boardID, ptttype.FN_DIR)
	// rename.Dir content
	err = cmsys.DeleteRecord(filePath, index, ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return err
	}
	return nil
}
