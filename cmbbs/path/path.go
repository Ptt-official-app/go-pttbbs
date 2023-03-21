package path

import (
	"os"
	"strconv"
	"strings"

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
	if !userID.IsValid() {
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

func SetBPath(boardID *ptttype.BoardID_t) string {
	return strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_BOARD,
		string(boardID[0]),
		types.CstrToString(boardID[:]),
	},
		string(os.PathSeparator),
	)
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

func SetBNFile(boardID *ptttype.BoardID_t, filename string, idx int) (string, error) {
	if filename[0] == '\x00' || !IsValidFilename(filename) {
		return "", ptttype.ErrInvalidFilename
	}

	theFilename := filename + "." + strconv.Itoa(idx)

	return strings.Join([]string{
		ptttype.BBSHOME,
		ptttype.DIR_BOARD,
		string(boardID[0]),
		types.CstrToString(boardID[:]),
		theFilename,
	},
		string(os.PathSeparator),
	), nil
}

func SetBBSHomePath(filename string) string {
	return ptttype.SetBBSHomePath(filename)
}

// SetDIRPath
//
// dirFilename: boards/[boardID]/.DIR
// filename: the filename in [boardID]
func SetDIRPath(dirFilename string, filename string) string {
	idx := strings.LastIndex(dirFilename, string(os.PathSeparator))
	theDir := dirFilename[:idx]
	return strings.Join([]string{
		theDir,
		filename,
	}, string(os.PathSeparator))
}
