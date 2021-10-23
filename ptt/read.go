package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func DeleteArticles(boardID *ptttype.BoardID_t, filename *ptttype.Filename_t, index ptttype.SortIdx) error {
	filePath, err := path.SetBFile(boardID, ptttype.FN_DIR)
	// rename.Dir content
	// index must transfer to sortIdxInStore
	err = cmsys.DeleteRecord(filePath, index.ToSortIdxInStore(), ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return err
	}
	return nil
}
