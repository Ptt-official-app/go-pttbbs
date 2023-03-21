package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetPostTemplate(uuserID UUserID, bboardID BBoardID, templateID ptttype.SortIdx, retrieveTS types.Time4, isHash bool) (content []byte, mtime types.Time4, sum cmsys.Fnv64_t, err error) {
	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, 0, ErrInvalidParams
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, 0, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, 0, 0, ErrInvalidParams
	}

	return ptt.ReadPostTemplate(userecRaw, uid, boardIDRaw, bid, templateID, retrieveTS, isHash)
}
