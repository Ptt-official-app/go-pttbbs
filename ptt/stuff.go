package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func is_uBM(userID *ptttype.UserID_t, bm *ptttype.BM_t) bool {
	userIDBytes := types.CstrToBytes(userID[:])
	bmBytes := types.CstrToBytes(bm[:])
	theIdx := types.Cstrstr(bmBytes, userIDBytes)
	if theIdx < 0 {
		return false
	}

	isValidHead := true
	if theIdx > 0 {
		isValidHead = !types.Isalnum(bmBytes[theIdx-1])
	}

	isValidTail := true
	if theIdx+len(userIDBytes) < len(bmBytes) {
		isValidTail = !types.Isalnum(bmBytes[theIdx+len(userIDBytes)])
	}

	return isValidHead && isValidTail
}

//setBDir
//
//XXX do not deal with digest-mode in go-pttbbs.
//    will use middleware for digest.
func setBDir(boardID *ptttype.BoardID_t) (filename string, err error) {
	return path.SetBFile(boardID, ptttype.FN_DIR)
}

//setBBottom
//
//XXX do not deal with digest-mode in go-pttbbs.
//    will use middleware for digest.
func setBBottom(boardID *ptttype.BoardID_t) (filename string, err error) {
	return path.SetBFile(boardID, ptttype.FN_DIR_BOTTOM)
}
