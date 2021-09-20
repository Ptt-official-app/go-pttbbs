package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
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
	posttype []byte,
	title []byte,
	content [][]byte,
	oldSZ int,
	oldsum cmsys.Fnv64_t,
	ip string) (newContent []byte, mtime types.Time4, newTitle []byte, newClass []byte, err error) {

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	userIDRaw, err := uuserID.ToRaw()
	if err != nil {
		return nil, 0, nil, nil, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	bid, boardIDRaw, err := bboardID.ToRaw()
	if err != nil {
		return nil, 0, nil, nil, ErrInvalidParams
	}

	from := fromd.GetFrom(ipRaw)

	filename := articleID.ToFilename()
	newContent, mtime, newTitleRaw, err := ptt.EditPost(userecRaw, uid, boardIDRaw, bid, filename, posttype, title, content, oldSZ, oldsum, ipRaw, from)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	_, realTitleWithClass := cmbbs.SubjectEx(newTitleRaw)

	newClass, newTitle = titleToClass(realTitleWithClass)

	return newContent, mtime, newTitle, newClass, nil
}
