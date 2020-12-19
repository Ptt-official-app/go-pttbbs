package bbs

import (
	"strconv"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BBoardID string //The integrated bid-boardID, concat with _, safe because bid is number >= 1.

//ToBBoardID
//
//concat bid and boardID as BBoardID, concat with _
//because bid and boardIDRaw are from ptt, no need to check the validity
func ToBBoardID(bid ptttype.Bid, boardIDRaw *ptttype.BoardID_t) BBoardID {
	if !bid.IsValid() {
		return BBoardID("")
	}
	if !boardIDRaw.IsValid() {
		return BBoardID("")
	}
	return BBoardID(bid.String() + "_" + types.CstrToString(boardIDRaw[:]))
}

//ToRaw
//
//BBoardID is possible coming from outside, requiring validation.
func (b BBoardID) ToRaw() (bid ptttype.Bid, boardIDRaw *ptttype.BoardID_t, err error) {
	bList := strings.Split(string(b), "_")
	if len(bList) != 2 {
		return 0, nil, ErrInvalidBBoardID
	}

	// check bid validity
	bid_i, err := strconv.Atoi(bList[0])
	if err != nil {
		return 0, nil, err
	}
	bid = ptttype.Bid(bid_i)
	if !bid.IsValid() {
		return 0, nil, ErrInvalidBBoardID
	}

	//boardIDRaw
	boardIDRaw = &ptttype.BoardID_t{}
	copy(boardIDRaw[:], []byte(bList[1]))
	if !boardIDRaw.IsValid() {
		return 0, nil, ErrInvalidBBoardID
	}

	return bid, boardIDRaw, nil
}

func (b BBoardID) ToBrdname() string {
	_, boardIDRaw, _ := b.ToRaw()
	return types.CstrToString(boardIDRaw[:])
}
