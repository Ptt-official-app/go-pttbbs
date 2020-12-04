package cmsys

import "os"

func GetNumRecords(filename string, size uintptr) int {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0
	}

	return int(stat.Size() / int64(size))
}
