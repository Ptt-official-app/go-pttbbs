package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

//boardPermStat
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L197
//
//The original hasBoardPerm
func boardPermStat(user *ptttype.UserecRaw, uid ptttype.Uid, board *ptttype.BoardHeaderRaw, bid ptttype.Bid) ptttype.BoardStatAttr {

	//SYSOP
	logrus.Infof("userLevel: %v SYSOP: %v", user.UserLevel, user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP))
	if user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return ptttype.NBRD_FAV
	}

	return boardPermStatNormally(user, uid, board, bid)
}

//boardPermStatNormally
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L157
//
//The original hasBoardPermNormally
//
//Original code was mixing-up BoardStat with this function, making the code not easy to comprehend.
//BM / Police / SYSOP treat the board as NBRD_FAV, while others treat the board as NBRD_BOARD.
//It's because that newBoardStat is hacked to forcely add BRD_POSTMASK if not set properly and the type is NBRD_BOARD.
//Need to figure out a better method to solve this issue.
func boardPermStatNormally(user *ptttype.UserecRaw, uid ptttype.Uid, board *ptttype.BoardHeaderRaw, bid ptttype.Bid) ptttype.BoardStatAttr {
	level := board.Level
	brdAttr := board.BrdAttr

	// allow POLICE to enter BM boards
	if (level&ptttype.PERM_BM != 0) && (user.UserLevel.HasUserPerm(ptttype.PERM_POLICE) || user.UserLevel.HasUserPerm(ptttype.PERM_POLICE_MAN)) {
		return ptttype.NBRD_FAV
	}

	/* 板主 */
	if IsBMCache(user, uid, bid) {
		return ptttype.NBRD_FAV
	}

	/* 祕密看板：核對首席板主的好友名單 */
	if brdAttr&ptttype.BRD_HIDE != 0 {
		bidInCache := bid.ToBidInStore()
		uidInCache := uid.ToUidInStore()
		if !cache.IsHiddenBoardFriend(bidInCache, uidInCache) {
			if brdAttr&ptttype.BRD_POSTMASK != 0 {
				return ptttype.NBRD_INVALID
			} else {
				//XXX return 2;
				//    what's this? (in addnewbrdstat, to set brd_postmask)
				//    need to simplify this function.
				return ptttype.NBRD_BOARD
			}
		} else {
			return ptttype.NBRD_FAV
		}
	}

	// TODO Change this to a query on demand.
	/* 十八禁看板 */
	if brdAttr&ptttype.BRD_OVER18 != 0 && !user.Over18 {
		return ptttype.NBRD_INVALID
	}

	if level != 0 && (brdAttr&ptttype.BRD_POSTMASK) == 0 && !user.UserLevel.HasUserPerm(level) {
		return ptttype.NBRD_INVALID
	}

	return ptttype.NBRD_FAV
}
