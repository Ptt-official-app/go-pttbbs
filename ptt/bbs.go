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
)

//ReadPost
//
//pmore is replaced by frontend.
//We just need to return the whole content.
//We do not update brc here, because it requires lots of file-disk.
//require middlewares to handle user-read-article.
func ReadPost(
	user *ptttype.UserecRaw,
	uid ptttype.Uid,
	boardID *ptttype.BoardID_t,
	bid ptttype.Bid,
	filename *ptttype.Filename_t,
	retrieveTS types.Time4) (

	content []byte,
	mtime types.Time4,
	err error,
) {

	//1. check valid filename
	if filename[0] == 'L' || filename[0] == 0 {
		return nil, 0, ErrInvalidParams
	}

	cache.StatInc(ptttype.STAT_READPOST)

	//2. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, 0, ErrNotPermitted
	}

	//3. get filename
	theFilename, err := path.SetBFile(boardID, filename.String())
	if err != nil {
		return nil, 0, err
	}

	//4. check mtime
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

	//XXX do not do brc for now.
	//brcAddList(boardID, filename, updateTS)

	return content, mtime, nil
}

func getBoardRestrictionReason(
	user *ptttype.UserecRaw,
	uid ptttype.Uid,
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
	uid ptttype.Uid,
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
	uid ptttype.Uid,
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

	//check-post-perm2
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

	//do not permit user without loginok.
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

	if isUseAnony {
		postfile.Filemode |= ptttype.FILE_ANONYMOUS
		postfile.SetAnonUID(uid)
	} else {
		postfile.Modified = types.DashT(fpath)
		postfile.SetMoney(0)
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
		cmsys.LogFilef(fpath, cmsys.LOG_CREAT, msg)
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

	//XXX no USE_POSTD for now.

	summary = ptttype.NewArticleSummaryRaw(idx, boardID, postfile)

	return summary, nil
}

func checkBoardAnonymous(board *ptttype.BoardHeaderRaw) bool {
	if !ptttype.HAVE_ANONYMOUS {
		return false
	}

	return board.BrdAttr&ptttype.BRD_ANONYMOUS != 0
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

func checkCooldown(user *ptttype.UserecRaw, uid ptttype.Uid, board *ptttype.BoardHeaderRaw, bid ptttype.Bid) (isCooldown bool, err error) {

	limit := [8]int{4000, 1, 2000, 2, 1000, 3, -1, 10}

	nowTS := types.NowTS()
	cooldownTime := cache.CooldownTimeOf(uid)
	diff := cooldownTime - nowTS
	if diff < 0 {
		cooldownTime &= 0x7FFFFFF0
		_ = cache.SetCooldownTime(uid, cooldownTime)

		return false, nil
	}

	//XXX ignoring currmode == MODE_BOARD
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
	uid ptttype.Uid,
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
	uid ptttype.Uid,
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

//isModeBoard
//https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L376
func isModeBoard(
	user *ptttype.UserecRaw,
	uid ptttype.Uid,
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
	uid ptttype.Uid,
	board *ptttype.BoardHeaderRaw,
	bid ptttype.Bid,
	ip *ptttype.IPv4_t,
	from []byte) (err error) {

	_, err = WriteFile(fpath, flags, true, isUseAnony, title, content, user, uid, board, bid, ip, from)

	return err
}

func doCrosspost(
	user *ptttype.UserecRaw,
	uid ptttype.Uid,
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

	//new title
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

	//filename
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

	//use filename.CreateTime as lastPosttime.
	lastPosttime, err := fh.Filename.CreateTime()
	if err != nil {
		return nil, err
	}
	_ = cache.SetLastPosttime(bid, lastPosttime)

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

func CheckPostPerm2(uid ptttype.Uid, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return CheckModifyPerm(uid, user, bid, board)
}

func CheckModifyPerm(uid ptttype.Uid, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return postpermMsg(uid, user, bid, board)
}
