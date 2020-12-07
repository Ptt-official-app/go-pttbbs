package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptt"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	log "github.com/sirupsen/logrus"
)

func LoadGeneralBoards(userID string, startIdx int32, nBoards int, keyword []byte) (summary []*BoardSummary, nextIdx int32, err error) {

	userIDRaw := &ptttype.UserID_t{}
	copy(userIDRaw[:], []byte(userID))

	uid, userecRaw, err := ptt.InitCurrentUser(userIDRaw)
	log.Infof("bbs.LoadGeneralBoards: after InitCurrentUser: uid: %v userecRaw: %v e: %v", uid, userecRaw, err)
	if err != nil {
		return nil, -1, err
	}

	log.Infof("bbs.LoadGeneralBoards: to ptt.LoadGeneralBoards: uid: %v userecRaw: %v startIdx: %v nBoards: %v keyword: %v", uid, userecRaw, startIdx, nBoards, keyword)
	summaryRaw, nextIdx, err := ptt.LoadGeneralBoards(userecRaw, uid, startIdx, nBoards, keyword)
	if err != nil {
		return nil, -1, err
	}

	summary = make([]*BoardSummary, len(summaryRaw))
	log.Infof("bbs.LoadGeneralBoards: summaryRaw: %v", len(summaryRaw))
	for idx, each := range summaryRaw {
		eachSummary := NewBoardSummaryFromRaw(each)
		summary[idx] = eachSummary
	}

	return summary, nextIdx, nil
}
