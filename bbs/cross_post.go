package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/fromd"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func CrossPost(
	uuserID UUserID,
	bboardID BBoardID,
	articleID ArticleID,
	xboardID BBoardID,
	ip string,
) (articleSummary *ArticleSummary, comment []byte, commentMTime types.Time4, err error) {

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, nil, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, nil, 0, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, nil, 0, ErrInvalidParams
	}

	xbid, xboardIDRaw, err := xboardID.ToRaw()
	if err != nil {
		return nil, nil, 0, ErrInvalidParams
	}

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))
	from := fromd.GetFrom(ipRaw)

	filename := articleID.ToFilename()

	articleSummaryRaw, comment, commentMTime, err := ptt.CrossPost(userecRaw, uid, boardIDRaw, bid, filename, xboardIDRaw, xbid, ptttype.FILE_NONE, ipRaw, from)
	if err != nil {
		return nil, nil, 0, err
	}

	articleSummary = NewArticleSummaryFromRaw(xboardID, articleSummaryRaw)

	return articleSummary, comment, commentMTime, err
}
