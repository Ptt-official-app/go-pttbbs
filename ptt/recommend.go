package ptt

import (
	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

//https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L2768

func Recommend(
	user *ptttype.UserecRaw,
	uid ptttype.Uid,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	filename *ptttype.Filename_t,

	commentType ptttype.CommentType,
	content []byte,
	ip *ptttype.IPv4_t,
	from []byte) (comment []byte, mtime types.Time4, err error) {
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, 0, ErrNotPermitted
	}

	//check-post-perm2
	err = CheckPostPerm2(uid, user, bid, board)
	if err != nil {
		return nil, 0, err
	}

	reason, err := getBoardRestrictionReason(user, uid, board, bid)
	if err != nil {
		return nil, 0, err
	}
	if reason != ptttype.RESTRICT_REASON_NONE {
		return nil, 0, ErrNotPermitted
	}

	isCooldown, err := checkCooldown(user, uid, board, bid)
	if err != nil {
		return nil, 0, err
	}
	if isCooldown {
		return nil, 0, ErrCooldown
	}

	dirFilename, err := setBDir(boardID)
	if err != nil {
		return nil, 0, err
	}
	logrus.Infof("Recommmend: boardID: %v dirFilename: %v", string(boardID[:]), dirFilename)

	//check record permission
	total, err := cache.GetBTotalWithRetry(bid)
	logrus.Infof("Recommmend: after GetGetBTotalWithRetry: total: %v e: %v", total, err)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, ErrInvalidParams
	}

	idx, fhdr, err := cmsys.GetRecord(dirFilename, filename, int(total))
	logrus.Infof("Recommmend: after GetRecord: idx: %v fhdr: %v e: %v", idx, fhdr, err)
	if err != nil {
		return nil, 0, err
	}
	if (board.BrdAttr&ptttype.BRD_NORECOMMEND) != 0 || filename[0] == 'L' || ((fhdr.Filemode&ptttype.FILE_MARKED) != 0 && (fhdr.Filemode&ptttype.FILE_SOLVED) != 0) {
		return nil, 0, ErrNotPermitted
	}

	//put recommand
	comment, err = FormatCommentString(user, board, commentType, content, ip, from)
	if err != nil {
		return nil, 0, err
	}

	mtime, err = doAddRecommend(dirFilename, idx, fhdr, comment, commentType)
	if err != nil {
		return nil, 0, err
	}

	return comment, mtime, nil

}
