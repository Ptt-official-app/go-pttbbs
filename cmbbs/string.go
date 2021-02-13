package cmbbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func SubjectEx(title *ptttype.Title_t) (titleType ptttype.SubjectType, newTitle []byte) {
	pTitle := types.CstrToBytes(title[:])

	titleType = ptttype.SUBJECT_NORMAL
	for {
		if len(pTitle) == 0 {
			break
		}
		if cmsys.StrcaseStartsWith(pTitle, ptttype.STR_REPLY) {
			pTitle = pTitle[len(ptttype.STR_REPLY):]
			titleType = ptttype.SUBJECT_REPLY
		} else if cmsys.StrcaseStartsWith(pTitle, ptttype.STR_FORWARD) {
			pTitle = pTitle[len(ptttype.STR_FORWARD):]
			titleType = ptttype.SUBJECT_FORWARD
		} else if ptttype.USE_LEGACY_FORWARD && cmsys.StrcaseStartsWith(pTitle, ptttype.STR_LEGACY_FORWARD) {
			pTitle = pTitle[len(ptttype.STR_LEGACY_FORWARD):]
			titleType = ptttype.SUBJECT_FORWARD
		} else {
			break
		}

		if len(pTitle) == 0 {
			break
		}
		if pTitle[0] == ' ' {
			pTitle = pTitle[1:]
		}
	}

	return titleType, pTitle
}
