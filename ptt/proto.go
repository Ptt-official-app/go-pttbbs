package ptt

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

func isVisible(meID ptttype.UtmpID, me *ptttype.UserInfoRaw, uentID ptttype.UtmpID, uentp *ptttype.UserInfoRaw) bool {
	return isVisibleStat(me, uentp, friendStat(meID, me, uentID, uentp))
}
