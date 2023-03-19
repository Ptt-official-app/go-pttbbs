package ptt

import (
	"bytes"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
)

func StripANSIMoveCmd(line []byte) (newLine []byte) {
	newLine = line
	for len(line) > 0 {
		idx := bytes.Index(newLine, []byte{ansi.ESC_CHR})
		if idx < 0 {
			break
		}
		newLine = newLine[idx:]
		if len(newLine) < 1 {
			break
		}
		newLine = newLine[1:]
		//nolint:revive
		for ; len(newLine) > 0 && bytes.IndexByte(ptttype.PATTERN_ANSI_CODE, newLine[0]) >= 0; newLine = newLine[1:] {
		}

		if len(newLine) == 0 {
			break
		}
		if bytes.IndexByte(ptttype.PATTERN_ANSI_MOVECMD, newLine[0]) >= 0 {
			newLine[0] = 's'
		}

	}

	return line
}
