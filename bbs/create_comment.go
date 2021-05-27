package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/fromd"
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func CreateComment(userID UUserID, boardID BBoardID, articleID ArticleID, commentType ptttype.CommentType, content []byte, ip string) (comment []byte, mtime types.Time4, err error) {

	ipRaw := &ptttype.IPv4_t{}
	copy(ipRaw[:], []byte(ip))

	userIDRaw, err := userID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return nil, 0, err
	}

	bid, boardIDRaw, err := boardID.ToRaw()
	if err != nil {
		return nil, 0, ErrInvalidParams
	}

	from := fromd.GetFrom(ipRaw)

	filename := articleID.ToFilename()

	logrus.Infof("bbs.CreateComment: to ptt.Recommend: commentType: %v", commentType)
	comment, mtime, err = ptt.Recommend(userecRaw, uid, boardIDRaw, bid, filename, commentType, content, ipRaw, from)
	if err != nil {
		return nil, 0, err
	}

	return comment, mtime, nil
}
