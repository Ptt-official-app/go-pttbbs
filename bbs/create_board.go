package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func CreateBoard(
	userID UUserID,
	clsBid ptttype.Bid,
	brdname string,
	brdClass []byte,
	brdTitle []byte,
	BMs []UUserID,
	brdAttr ptttype.BrdAttr,
	level ptttype.PERM,
	chessCountry ptttype.ChessCode,
	isGroup bool,
) (boardID BBoardID, err error) {

	userIDRaw, err := userID.ToRaw()
	if err != nil {
		return "", ErrInvalidParams
	}

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	if err != nil {
		return "", err
	}

	brdnameRaw := &ptttype.BoardID_t{}
	copy(brdnameRaw[:], []byte(brdname))

	BMs_userIDRaw := make([]*ptttype.UserID_t, len(BMs))
	for idx, each := range BMs {
		copy(BMs_userIDRaw[idx][:], []byte(each))
	}
	bms := ptttype.NewBM(BMs_userIDRaw)

	newBoardRaw, newBid, err := ptt.NewBoard(
		userecRaw,
		uid,
		clsBid,
		brdnameRaw,
		brdClass,
		brdTitle,
		bms,
		brdAttr,
		level,
		chessCountry,
		isGroup,
	)
	if err != nil {
		return "", err
	}

	boardID = ToBBoardID(newBid, &newBoardRaw.Brdname)

	return boardID, nil
}
