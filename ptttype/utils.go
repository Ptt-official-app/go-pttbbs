package ptttype

import "os"

func ValidUSHMEntry(x UtmpID) bool {
	return x >= 0 && x < USHM_SIZE
}

func SetBBSHomePath(filename string) string {
	return BBSHOME + string(os.PathSeparator) + filename
}
