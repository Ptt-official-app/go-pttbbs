package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/fromd"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func CreateArticle(
	uuserID UUserID,
	bboardID BBoardID,
	posttype []byte,
	title []byte,
	content [][]byte,
	ip string) (

	summary *ArticleSummary,
	err error) {

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, ErrInvalidParams
	}

	from := fromd.GetFrom(ipRaw)

	summaryRaw, err := ptt.NewPost(userecRaw, uid, boardIDRaw, bid, posttype, title, content, ipRaw, from)
	if err != nil {
		return nil, err
	}

	summary = NewArticleSummaryFromRaw(bboardID, summaryRaw)

	return summary, nil
}
