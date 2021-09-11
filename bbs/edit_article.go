package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/fromd"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func EditArticle(
	uuserID UUserID,
	bboardID BBoardID,
	articleID ArticleID,
	title []byte,
	content [][]byte,
	oldSZ int,
	oldsum cmsys.Fnv64_t,
	ip string) (newContent []byte, mtime types.Time4, err error) {

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, 0, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	from := fromd.GetFrom(ipRaw)

	filename := articleID.ToFilename()
	return ptt.EditPost(userecRaw, uid, boardIDRaw, bid, filename, title, content, oldSZ, oldsum, ipRaw, from)
}
