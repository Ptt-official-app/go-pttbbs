package cmsys

import (
	"bufio"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func FileExistsRecord(filename string, key []byte) bool {
	return FileFindRecord(filename, key) > 0
}

//FileFindRecord (starting from 1)
func FileFindRecord(filename string, key []byte) (idx int) {
	file, err := os.Open(filename)
	if err != nil {
		return 0
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for line, err := types.ReadLine(reader); err == nil; line, err = types.ReadLine(reader) {

		idx++

		first, _ := tokenize(line, ptttype.BYTES_SPACE)
		if types.Cstrcasecmp(key, first) == 0 {
			return idx
		}
	}

	return 0
}

func tokenize(line []byte, sep []byte) (first []byte, theRest []byte) {
	endIdx := len(line)
	for idx, each := range line {
		for _, each2 := range sep {
			if each == each2 {
				endIdx = idx
				break
			}
		}
	}
	return line[:endIdx], line[endIdx:]
}
