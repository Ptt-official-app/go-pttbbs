package ptt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
)

//ReadPost
//
//pmore is replaced by frontend.
//We just need to return the whole content.
//We do not update brc here, because it requires lots of file-disk.
//require middlewares to handle user-read-article.
func ReadPost(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	filename *ptttype.Filename_t,
	retrieveTS types.Time4) (

	content []byte,
	mtime types.Time4,
	err error,
) {
	// 1. check valid filename
	if filename[0] == 'L' || filename[0] == 0 {
		return nil, 0, ErrInvalidParams
	}

	_ = cache.StatInc(ptttype.STAT_READPOST)

	// 2. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, 0, ErrNotPermitted
	}

	// 3. get filename
	theFilename, err := path.SetBFile(boardID, filename.String())
	if err != nil {
		return nil, 0, err
	}

	// 4. check mtime
	stat, err := os.Stat(theFilename)
	if err != nil {
		return nil, 0, err
	}
	mtime = types.TimeToTime4(stat.ModTime())
	if mtime <= retrieveTS {
		return nil, mtime, nil
	}

	file, err := os.Open(theFilename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	content, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	// XXX do not do brc for now.
	// brcAddList(boardID, filename, updateTS)

	return content, mtime, nil
}

// CheckPostRestriction
//
// true if pass the check (can post)
func CheckPostRestriction(user *ptttype.UserecRaw, uid ptttype.UID, board *ptttype.BoardHeaderRaw, bid ptttype.Bid) bool {
	reason, err := getBoardRestrictionReason(user, uid, board, bid)
	return err == nil && reason == ptttype.RESTRICT_REASON_NONE
}

func getBoardRestrictionReason(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid) (

	reason ptttype.RestrictReason,
	err error) {
	if user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return ptttype.RESTRICT_REASON_NONE, nil
	}

	if IsBMCache(user, uid, bid) {
		return ptttype.RESTRICT_REASON_NONE, nil
	}

	return getRestrictionReason(user.NumLoginDays, user.BadPost, board.PostLimitLogins, board.PostLimitBadpost)
}

//NewPost
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L1624
func NewPost(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	posttype []byte,
	title []byte,
	content [][]byte,
	ip *ptttype.IPv4_t,
	from []byte) (

	summary *ptttype.ArticleSummaryRaw,
	err error) {
	return DoPostArticle(user, uid, boardID, bid, ptttype.EDITFLAG_KIND_NEWPOST, posttype, title, content, ip, from)
}

//DoPostArticle
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L1258
func DoPostArticle(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	flags ptttype.EditFlag,
	posttype []byte,
	title []byte,
	content [][]byte,
	ip *ptttype.IPv4_t,
	from []byte) (

	summary *ptttype.ArticleSummaryRaw,
	err error) {

	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, ErrNotPermitted
	}

	// check-post-perm2
	err = CheckPostPerm2(uid, user, bid, board)
	if err != nil {
		return nil, err
	}

	reason, err := getBoardRestrictionReason(user, uid, board, bid)
	if err != nil {
		return nil, err
	}
	if reason != ptttype.RESTRICT_REASON_NONE {
		return nil, ErrNotPermitted
	}

	isCooldown, err := checkCooldown(user, uid, board, bid)
	if err != nil {
		return nil, err
	}
	if isCooldown {
		return nil, ErrCooldown
	}

	// do not permit user without loginok.
	if !user.UserLevel.HasUserPerm(ptttype.PERM_LOGINOK) {
		return nil, ErrNotPermitted
	}

	_ = cache.SetUtmpMode(uid, ptttype.USER_OP_POSTING)

	postfile := &ptttype.FileHeaderRaw{}

	boardPath := path.SetBPath(boardID)
	fpath, err := cmbbs.Stampfile(boardPath, postfile)
	if err != nil {
		return nil, err
	}

	fullTitle := doPostArticleFullTitle(posttype, title)
	fullTitle = tnSafeStrip(fullTitle, user, uid, board, bid)

	isUseAnony := checkBoardAnonymous(board)
	ownerID := &user.UserID

	if isUseAnony {
		ownerID = ptttype.ANONYMOUS_ID
	}

	err = doPostArticleWriteFile(fpath, flags, isUseAnony, fullTitle, content, user, uid, board, bid, ip, from)
	if err != nil {
		return nil, err
	}

	if isUseAnony { // skip errcheck because what we focus is publishing the articles
		postfile.Filemode |= ptttype.FILE_ANONYMOUS
		_ = postfile.SetAnonUID(uid)
	} else {
		postfile.Modified = types.DashT(fpath)
		_ = postfile.SetMoney(0)
	}
	copy(postfile.Owner[:], ownerID[:])
	copy(postfile.Title[:], fullTitle)

	// Ptt: stamp file again to make it order
	//      fix the bug that search failure in getindex
	//      stampfile_u is used when you don't want to clear other fields
	fpath2, err := cmbbs.StampfileU(boardPath, postfile)
	if err != nil {
		return nil, err
	}

	if ptttype.QUERY_ARTICLE_URL {
		url := GetWebURL(board, postfile)
		msg := fmt.Sprintf("%s %v\n", ptttype.STR_URL_DISPLAYNAME_BIG5, url)
		_ = cmsys.LogFilef(fpath, cmsys.LOG_CREAT, msg)
	}

	bdir, err := setBDir(boardID)
	if err != nil {
		return nil, err
	}
	idx, err := cmsys.AppendRecord(bdir, postfile, ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return nil, err
	}

	err = os.Rename(fpath, fpath2)
	if err != nil {
		return nil, err
	}

	err = cache.SetBTotal(bid)
	if err != nil {
		return nil, err
	}

	_ = brcAddList(user, uid, board, bid, &postfile.Filename, postfile.Modified)

	if board.IsOpenBRD() {
		if user.NumLoginDays < ptttype.NEWIDPOST_LIMIT_DAYS {
			_, _ = doCrosspost(user, uid, board, bid, ptttype.BN_NEWIDPOST, postfile, fpath2)
		}

		if board.BrdAttr&ptttype.BRD_HIDE != 0 {
			_, _ = doCrosspost(user, uid, board, bid, ptttype.BN_ALLHIDPOST, postfile, fpath2)
		} else {
			_, _ = doCrosspost(user, uid, board, bid, ptttype.BN_ALLPOST, postfile, fpath2)
		}
	}

	if !isUseAnony {
		_ = pwcuIncNumPost(user, uid)
	}

	if board.BrdAttr&ptttype.BRD_ANONYMOUS != 0 {
		_, _ = doCrosspost(user, uid, board, bid, ptttype.BN_UNANONYMOUS, postfile, fpath2)
	}

	if ptttype.USE_COOLDOWN {
		nowTS := types.NowTS()
		if board.NUser > 30 {
			if cache.CooldownTimeOf(uid) < nowTS {
				_ = cache.AddCooldownTime(uid, 5)
			}
		}
		_ = cache.AddPosttimes(uid, 1)
	}

	// XXX no USE_POSTD for now.

	summary = ptttype.NewArticleSummaryRaw(idx, boardID, postfile)

	return summary, nil
}

func CrossPost(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	filename *ptttype.Filename_t,
	xBoardID *ptttype.BoardID_t,
	xBid ptttype.Bid,
	filemode ptttype.FileMode,
	ip *ptttype.IPv4_t,
	from []byte) (

	articleSummary *ptttype.ArticleSummaryRaw,
	comment []byte,
	commentMTime types.Time4,
	err error) {

	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, nil, 0, err
	}

	if board.BrdAttr.HasPerm(ptttype.BRD_VOTEBOARD) {
		return nil, nil, 0, ErrVoteBoard
	}

	// require board permitting user.
	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, nil, 0, ErrNotPermitted
	}

	// get fileheader
	fileIdx, fileHeader, err := getFileHeader(boardID, bid, filename)
	if err != nil {
		return nil, nil, 0, ptttype.ErrInvalidFilename
	}
	if fileHeader.Owner[0] == '-' {
		return nil, nil, 0, ErrDeleted
	}

	// 3. get filename
	theFilename, err := path.SetBFile(boardID, filename.String())
	if err != nil {
		return nil, nil, 0, err
	}

	// 4. check mtime
	_, err = os.Stat(theFilename)
	if err != nil {
		return nil, nil, 0, err
	}

	// XXX violate law
	if user.UserLevel.HasUserPerm(ptttype.PERM_VIOLATELAW) {
		return nil, nil, 0, ErrViolateLaw
	}

	if !user.UserLevel.HasUserPerm(ptttype.PERM_LOGINOK) {
		return nil, nil, 0, ErrNotLoginOk
	}

	if board.BrdAttr.HasPerm(ptttype.BRD_CPLOG) && (CheckPostPerm2(uid, user, bid, board) != nil || !CheckPostRestriction(user, uid, board, bid)) {
		return nil, nil, 0, ErrPermitNoPost
	}

	// XXX LOCAL_ALERT_CROSSPOST

	// check xboard
	xBoard, err := cache.GetBCache(xBid)
	if err != nil {
		return nil, nil, 0, err
	}

	if !hasPostPerm(user, uid, xBoard, xBid) {
		return nil, nil, 0, ErrPermitNoPost
	}

	reason, err := getBoardRestrictionReason(user, uid, xBoard, xBid)
	if err != nil || reason != ptttype.RESTRICT_REASON_NONE {
		return nil, nil, 0, ErrPermitNoPost
	}

	isCooldown, err := checkCooldown(user, uid, xBoard, xBid)
	if err != nil || isCooldown {
		return nil, nil, 0, ErrCooldown
	}

	// XXX crosspost verify captcha

	// do cross-post
	xBPath := path.SetBPath(xBoardID)
	xFileHeader := &ptttype.FileHeaderRaw{}
	xFilename, err := cmbbs.Stampfile(xBPath, xFileHeader)
	if err != nil {
		return nil, nil, 0, err
	}
	copy(xFileHeader.Owner[:], user.UserID[:])
	title := bytes.Join([][]byte{ptttype.STR_FORWARD, types.CstrToBytes(fileHeader.Title[:])}, []byte{' '})
	copy(xFileHeader.Title[:], title)
	types.TrimDBCS(xFileHeader.Title[:])

	xFileHeader.Filemode = filemode

	fname, err := path.SetBFile(boardID, fileHeader.Filename.String())
	if err != nil {
		return nil, nil, 0, err
	}

	// write-file
	err = crossPostWriteFile(xFilename, types.CstrToBytes(xFileHeader.Title[:]), user, xBoard, fname, board, boardID, filename, from, ip)
	if err != nil {
		return nil, nil, 0, err
	}

	// log in allpost
	err = logCrosspostInAllpost(user, xBoard, xBid, xFileHeader)
	if err != nil {
		return nil, nil, 0, err
	}

	// add recommend to original board
	comment, err = crossPostComment(user, uid, board, bid, xBoard, xBid, xBoardID, from, ip)
	if err != nil {
		return nil, nil, 0, err
	}

	bdir, err := setBDir(boardID)
	if err != nil {
		return nil, nil, 0, err
	}
	if comment != nil {
		commentMTime, err = doAddRecommend(bdir, fileIdx, fileHeader, comment, ptttype.COMMENT_TYPE_FORWARD)
		if err != nil {
			return nil, nil, 0, err
		}
	}

	xbdir, err := setBDir(xBoardID)
	if err != nil {
		return nil, nil, 0, err
	}
	idx, err := cmsys.AppendRecord(xbdir, xFileHeader, ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return nil, nil, 0, err
	}

	if ptttype.USE_COOLDOWN {
		if xBoard.NUser > 30 {
			nowTS := types.NowTS()
			if cache.CooldownTimeOf(uid) < nowTS {
				_ = cache.AddCooldownTime(uid, 5)
			}
		}
		_ = cache.AddPosttimes(uid, 1)
	}

	_ = cache.SetBTotal(xBid)

	articleSummary = ptttype.NewArticleSummaryRaw(idx, xBoardID, xFileHeader)

	return articleSummary, comment, commentMTime, nil
}

func crossPostWriteFile(xFilename string, title []byte, user *ptttype.UserecRaw, xboard *ptttype.BoardHeaderRaw, fname string, board *ptttype.BoardHeaderRaw, boardID *ptttype.BoardID_t, filename *ptttype.Filename_t, from []byte, ip *ptttype.IPv4_t) (err error) {
	file, err := os.OpenFile(xFilename, os.O_CREATE|os.O_WRONLY, ptttype.DEFAULT_FILE_CREATE_PERM)
	if err != nil {
		return err
	}
	defer file.Close()

	err = writeHeader(file, ptttype.EDITFLAG_NONE, title, user, xboard)
	if err != nil {
		return err
	}

	if !board.IsOpenBRD() {
		_, err = file.Write(CROSS_POST_HIDDEN_BOARD)
		if err != nil {
			return err
		}
		err = bSuckinfileInvis(file, fname, board)
		if err != nil {
			return err
		}
	} else {
		aidc := filename.ToAidu().ToAidc()
		theBytes := bytes.Join([][]byte{CROSS_POST_BOARD_PREFIX, types.CstrToBytes(boardID[:]), CROSS_POST_BOARD_INFIX, types.CstrToBytes(aidc[:]), CROSS_POST_BOARD_POSTFIX}, []byte{})

		_, err = file.Write(theBytes)
		if err != nil {
			return err
		}

		err = bSuckinfile(file, fname)
		if err != nil {
			return err
		}
	}

	isUseAnony := checkBoardAnonymous(xboard) || checkBoardAnonymous(board)
	err = addForwardSignature(file, user, isUseAnony, ip, from)
	if err != nil {
		return err
	}

	return nil
}

func logCrosspostInAllpost(user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw, bid ptttype.Bid, fileHeader *ptttype.FileHeaderRaw) (err error) {
	if !board.IsOpenBRD() {
		return nil
	}

	bidAllpost, err := cache.GetBid(ptttype.BN_ALLPOST)
	if err != nil {
		return err
	}

	fh := &ptttype.FileHeaderRaw{}
	*fh = *fileHeader
	fh.Filemode = ptttype.FILE_LOCAL
	nowTS := types.NowTS()
	fh.Modified = nowTS

	// strlcpy(fh.owner, cuser.userid, sizeof(fh.owner))
	fh.Owner = ptttype.Owner_t{}
	copy(fh.Owner[:], user.UserID[:])

	// fh.title
	brdnameBytes := types.CstrToBytes(board.Brdname[:])
	// '…' appears for t_columns-33.
	theLen := 42 - len(brdnameBytes)
	titleBytes := types.CstrToBytes(board.Title[:])
	if len(titleBytes) > theLen {
		titleBytes = titleBytes[:theLen]
		titleBytes[theLen-2] = '\xa1'
		titleBytes[theLen-1] = 'K'
	}

	fwTitleStr := fmt.Sprintf("%s %-*.*s(%s)", ptttype.STR_FORWARD, theLen, theLen, titleBytes, brdnameBytes)
	fh.Title = ptttype.Title_t{}
	copy(fh.Title[:], []byte(fwTitleStr))

	bdir, err := setBDir(ptttype.BN_ALLPOST)
	if err != nil {
		return err
	}

	_, err = cmsys.AppendRecord(bdir, fh, ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return err
	}

	_ = cache.SetLastPosttime(bidAllpost, nowTS)
	_ = cache.TouchBPostNum(bidAllpost, 1)

	return nil
}

func crossPostComment(user *ptttype.UserecRaw, uid ptttype.UID, board *ptttype.BoardHeaderRaw, bid ptttype.Bid, xBoard *ptttype.BoardHeaderRaw, xBid ptttype.Bid, xBoardID *ptttype.BoardID_t, from []byte, ip *ptttype.IPv4_t) (comment []byte, err error) {
	if !board.BrdAttr.HasPerm(ptttype.BRD_CPLOG) {
		return nil, nil
	}

	// XXX magic number
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L2243
	maxLength := 51 + 2 - 6
	bnameStr := ""
	if xBoard.IsOpenBRD() {
		bnameStr = fmt.Sprintf("%s %s", CROSS_POST_COMMENT_BOARD, types.CstrToString(xBoardID[:]))
	} else {
		bnameStr = string(CROSS_POST_COMMENT_HIDDEN_BOARD)
	}
	maxLength -= types.Cstrlen(user.UserID[:]) + len(bnameStr)

	tail := ""
	nowTS := types.NowTS()
	if ptttype.GUESTRECOMMEND {
		tail = fmt.Sprintf("%15s %s", types.CstrToString(ip[:]), nowTS.CdateMd())
	} else {
		// XXX magic number
		// https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L2262
		maxLength += (15 - 6)
		tail = fmt.Sprintf(" %s", nowTS.CdateMdHM())
	}

	commentStr := fmt.Sprintf(CROSS_POST_COMMENT_PREFIX+"%s"+CROSS_POST_COMMENT_INFIX+"%s"+ansi.ANSI_RESET_STR+"%*s%s\n", types.CstrToString(user.UserID[:]), bnameStr, maxLength, "", tail)

	return []byte(commentStr), nil
}

func checkBoardAnonymous(board *ptttype.BoardHeaderRaw) bool {
	if !ptttype.HAVE_ANONYMOUS {
		return false
	}

	return board.BrdAttr.HasPerm(ptttype.BRD_ANONYMOUS)
}

//doPostArticleFullTitle
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L1348
func doPostArticleFullTitle(posttype []byte, title []byte) (fullTitle []byte) {
	if len(posttype) == 0 {
		return title
	}

	return bytes.Join([][]byte{{'['}, posttype, {']', ' '}, title}, []byte{})
}

func checkCooldown(user *ptttype.UserecRaw, uid ptttype.UID, board *ptttype.BoardHeaderRaw, bid ptttype.Bid) (isCooldown bool, err error) {
	limit := [8]int{4000, 1, 2000, 2, 1000, 3, -1, 10}

	nowTS := types.NowTS()
	cooldownTime := cache.CooldownTimeOf(uid)
	diff := cooldownTime - nowTS
	if diff < 0 {
		cooldownTime &= 0x7FFFFFF0
		_ = cache.SetCooldownTime(uid, cooldownTime)

		return false, nil
	}

	// XXX ignoring currmode == MODE_BOARD
	if user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) {
		return false, nil
	}

	if board.BrdAttr&ptttype.BRD_COOLDOWN != 0 {
		return true, nil
	}

	posttimes := cache.PosttimesOf(uid)

	if posttimes == 0xf {
		return true, nil
	}

	if ptttype.REJECT_FLOOD_POST {
		for i := 0; i < len(limit); i += 2 {
			if int(board.NUser) > limit[i] && int(posttimes) >= limit[i+1] {
				return true, nil
			}
		}
	}

	return false, nil
}

func tnSafeStrip(
	title []byte,
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid) (

	newTitle []byte) {

	if isTnAllowed(title, user, uid, board, bid) {
		return title
	}

	title = title[len(ptttype.TN_ANNOUNCE_BIG5):]

	return title
}

func isTnAllowed(
	title []byte,
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid) (

	isAllowed bool) {
	if ptttype.ALLOW_FREE_TN_ANNOUNCE {
		return true
	}

	// TN_ANNOUNCE is prohibited for non-BMs
	if isModeBoard(user, uid, board, bid) || user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP|ptttype.PERM_ACCOUNTS|ptttype.PERM_BOARD|ptttype.PERM_BBSADM|ptttype.PERM_VIEWSYSOP|ptttype.PERM_POLICE_MAN) {
		return true
	}

	// Note: 關於 subgroup op 的判定目前也是一團糟 - 小組長要從自己的分類
	// 進去才會有 GROUPOP(). 不過反正小組長跟群組長的人沒那麼多，就開放他們
	// always 可以使用 TN_ANNOUNCE 吧。
	if user.UserLevel.HasUserPerm(ptttype.PERM_SYSSUPERSUBOP | ptttype.PERM_SYSSUBOP) {
		return true
	}

	return !isTnAnnounce(title)
}

func isTnAnnounce(title []byte) (isValid bool) {
	return bytes.Equal(title[:len(ptttype.TN_ANNOUNCE_BIG5)], ptttype.TN_ANNOUNCE_BIG5)
}

// isModeBoard
// https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L376
func isModeBoard(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid) bool {
	if user.UserLevel.HasUserPerm(ptttype.PERM_NOCITIZEN) {
		return false
	}

	if user.UserLevel.HasUserPerm(ptttype.PERM_BOARD) ||
		IsBMCache(user, uid, bid) ||
		board.BM[0] <= ' ' && groupOp(user, uid, board) {

		return true
	}

	return false
}

func doPostArticleWriteFile(
	fpath string,
	flags ptttype.EditFlag,
	isUseAnony bool,
	title []byte,
	content [][]byte,
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid,
	ip *ptttype.IPv4_t,
	from []byte) (err error) {

	_, err = WriteFile(fpath, flags, true, isUseAnony, title, content, user, uid, board, bid, ip, from)

	return err
}

func doCrosspost(
	user *ptttype.UserecRaw,
	uid ptttype.UID,
	origBoard *ptttype.BoardHeaderRaw,
	origBid ptttype.Bid,
	boardID *ptttype.BoardID_t,
	header *ptttype.FileHeaderRaw,
	fpath string) (

	summary *ptttype.ArticleSummaryRaw,
	err error) {

	bid, err := cache.GetBid(boardID)
	if err != nil {
		return nil, err
	}
	if !bid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	titleType, title := cmbbs.SubjectEx(&header.Title)

	var prefix []byte
	switch titleType {
	case ptttype.SUBJECT_REPLY:
		prefix = ptttype.STR_REPLY
	case ptttype.SUBJECT_FORWARD:
		prefix = ptttype.STR_FORWARD
	}

	fh := &ptttype.FileHeaderRaw{}
	*fh = *header

	filename, err := path.SetBFile(boardID, types.CstrToString(header.Filename[:]))
	if err != nil {
		return nil, err
	}
	if types.Cstrcasecmp(boardID[:], ptttype.BN_UNANONYMOUS[:]) == 0 {
		fh.Owner = ptttype.Owner_t{}
		copy(fh.Owner[:], user.UserID[:])
	}

	// new title
	newTitle := make([]byte, 0, ptttype.TTLEN)
	if prefix != nil {
		newTitle = append(newTitle, prefix...)
		newTitle = append(newTitle, ' ')
	}
	currboardBytes := types.CstrToBytes(origBoard.Brdname[:])
	theLen := 42 - len(currboardBytes)
	newTitle = append(newTitle, dbcsSafeTrimTitle(title, theLen)...)
	newTitle = append(newTitle, '(')
	newTitle = append(newTitle, currboardBytes...)
	newTitle = append(newTitle, ')')

	fh.Title = ptttype.Title_t{}
	copy(fh.Title[:], newTitle)

	// filename
	if types.DashS(filename) > 0 {
		nowTS := types.NowTS()
		theLog := fmt.Sprintf("%v %v->%v %v: %v\n", nowTS.Cdatelite(), string(currboardBytes), types.CstrToString(boardID[:]), types.CstrToString(fh.Filename[:]), types.CstrToString(fh.Title[:]))
		_ = cmsys.LogFilef("log/conflict.log", cmsys.LOG_CREAT, theLog)
	}

	err = types.CopyFileToFile(fpath, filename)
	if err != nil {
		return nil, err
	}

	nowTS := types.NowTS()
	fh.Filemode = ptttype.FILE_LOCAL
	fh.Modified = nowTS

	bdir, err := setBDir(boardID)
	if err != nil {
		return nil, err
	}
	idx, err := cmsys.AppendRecord(bdir, fh, ptttype.FILE_HEADER_RAW_SZ)
	if err != nil {
		return nil, err
	}

	// use filename.CreateTime as lastPosttime.
	_ = cache.SetLastPosttime(bid, nowTS)

	_ = cache.TouchBPostNum(bid, 1)

	summary = ptttype.NewArticleSummaryRaw(idx, boardID, fh)

	return summary, nil
}

func dbcsSafeTrimTitle(title []byte, theLen int) (newTitle []byte) {
	if len(title) <= theLen {
		return title
	}

	title = title[:theLen-2]
	title = cmsys.DBCSSafeTrim(title)
	title = append(title, ptttype.STR_DOTS...)

	return title
}

func GetWebURL(board *ptttype.BoardHeaderRaw, fhdr *ptttype.FileHeaderRaw) (url string) {
	folder := types.CstrToString(board.Brdname[:])
	fn := types.CstrToString(fhdr.Filename[:])
	ext := ".html"
	if ptttype.USE_AID_URL {
		aidc := fhdr.Filename.ToAidu().ToAidc()
		fn = types.CstrToString(aidc[:])
		ext = ""
	}

	return ptttype.URL_PREFIX + "/" + folder + "/" + fn + ext
}

func CheckPostPerm2(uid ptttype.UID, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return CheckModifyPerm(uid, user, bid, board)
}

func CheckModifyPerm(uid ptttype.UID, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return postpermMsg(uid, user, bid, board)
}

func getFileHeader(boardID *ptttype.BoardID_t, bid ptttype.Bid, filename *ptttype.Filename_t) (idx ptttype.SortIdx, fileHeader *ptttype.FileHeaderRaw, err error) {
	dirFilename, err := setBDir(boardID)
	if err != nil {
		return 0, nil, err
	}

	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, nil, ptttype.ErrInvalidFilename
	}

	idx, fileHeader, err = cmsys.GetRecord(dirFilename, filename, int(total))
	if err != nil {
		return 0, nil, err
	}

	return idx, fileHeader, nil
}
