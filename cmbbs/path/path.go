package path

import (
	"os"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/names"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func SetHomePath(userID *ptttype.UserID_t) string {
	return strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_HOME,
		string(userID[0]),
		types.CstrToString(userID[:]),
	},
		string(os.PathSeparator),
	)
}

func SetHomeFile(userID *ptttype.UserID_t, filename string) (string, error) {
	if !names.IsValidUserID(userID) {
		return "", ptttype.ErrInvalidUserID
	}
	if filename[0] == '\x00' || !IsValidFilename(filename) {
		return "", ptttype.ErrInvalidFilename
	}
	return strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_HOME,
		string(userID[0]),
		types.CstrToString(userID[:]),
		filename,
	},
		string(os.PathSeparator),
	), nil
}

func IsValidFilename(filename string) bool {
	return !strings.Contains(filename, "..")
}

func SetBFile(boardID *ptttype.BoardID_t, filename string) (string, error) {
	if filename[0] == '\x00' || !IsValidFilename(filename) {
		return "", ptttype.ErrInvalidFilename
	}

	return strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_BOARD,
		string(boardID[0]),
		types.CstrToString(boardID[:]),
		filename,
	},
		string(os.PathSeparator),
	), nil
}

func SetBBSHomePath(filename string) string {
	return ptttype.BBSHOME + string(os.PathSeparator) + filename
}
