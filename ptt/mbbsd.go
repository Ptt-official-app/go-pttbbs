package ptt

import (
	"bufio"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//LoginQuery
//
//Params
//	userID: userID
//	passwd: passwd
//	ip: ip
//
//Return
//	*UserecRaw: user
//  error: err
func LoginQuery(userID *ptttype.UserID_t, passwd []byte, ip *ptttype.IPv4_t) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {
	if !userID.IsValid() {
		//log.Errorf("LoginQuery: invalid user id: userID: %v", userID)
		return 0, nil, ptttype.ErrInvalidUserID
	}

	uid, user, err = InitCurrentUser(userID)
	if err != nil {
		return 0, nil, err
	}

	isValid, err := cmbbs.CheckPasswd(user.PasswdHash[:], passwd)
	if err != nil {
		cmbbs.LogAttempt(userID, ip, true)
		return 0, nil, err
	}

	if !isValid {
		cmbbs.LogAttempt(userID, ip, true)
		return 0, nil, ptttype.ErrInvalidUserID
	}

	return uid, user, nil
}

//Login
//
//adopted from the original start_client.
//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L1399
func Login(userID *ptttype.UserID_t, passwd []byte, ip *ptttype.IPv4_t) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {

	uid, user, err = LoginQuery(userID, passwd, ip)
	if err != nil {
		return 0, nil, err
	}

	//we don't do loadCurrentUser
	//because logattempt, ensure_user_agreement_version
	//should be in middleware.

	err = userLogin(uid, user, ip)
	if err != nil {
		return 0, nil, err
	}

	//XXX we don't do auto-close-polls here.
	//we should have another goroutine to do auto-close-polls.

	return uid, user, nil
}

func userLogin(uid ptttype.Uid, user *ptttype.UserecRaw, ip *ptttype.IPv4_t) (err error) {

	utmpID, uinfo, err := setupUtmp(uid, user, ip, ptttype.USER_OP_LOGIN)
	if err != nil {
		return err
	}

	//XXX We should have new stats for go-pttbbs
	//_ = cache.StatInc(ptttype.STAT_MBBSD_ENTER)

	//XXX skip mail-related for now.
	//currutmp->alerts |= load_mailalert(cuser.userid)

	//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L1219
	cache.Shm.CheckMaxUser()

	if !(user.UserLevel.HasUserPerm(ptttype.PERM_SYSOP) && user.UserLevel.HasUserPerm(ptttype.PERM_SYSOPHIDE)) {
		doAloha(utmpID, uinfo, ptttype.ALOHA_MSG)
	}

	return nil
}

//setupUtmp
//
//There will be only 1-login per user in this process.
func setupUtmp(uid ptttype.Uid, user *ptttype.UserecRaw, ip *ptttype.IPv4_t, op ptttype.UserOpMode) (utmpID ptttype.UtmpID, uinfo *ptttype.UserInfoRaw, err error) {
	if !ptttype.IS_UTMP {
		return 0, nil, nil
	}

	uinfo = newUserInfoRaw(uid, user, ip, op)

	utmpID, err = getNewUtmpEnt(uinfo)
	if err != nil {
		return 0, nil, err
	}

	return utmpID, uinfo, nil
}

//newUserInfoRaw
//
//XXX we need cmsys.StripNoneBig5,
//with which newUserInfoRaw cannot be in ptttype.
func newUserInfoRaw(uid ptttype.Uid, user *ptttype.UserecRaw, ip *ptttype.IPv4_t, op ptttype.UserOpMode) *ptttype.UserInfoRaw {

	fromIP := types.InetAddr(types.CstrToString(ip[:]))
	nowTS := types.NowTS()

	//XXX we can do stringNoneBig5 here.
	//because now:
	//1. it's http-session based connection.
	//2. the pid is fixed by user.
	//3. user.nickname should not be affected.
	uinfo := &ptttype.UserInfoRaw{
		Pid:      types.Pid_t(uid) + types.Pid_t(types.DEFAULT_PID_MAX),
		Uid:      uid,
		Mode:     op,
		UserID:   user.UserID,
		Nickname: user.Nickname,

		UserLevel: user.UserLevel,
		LastAct:   nowTS,
		FromIp:    fromIP,
	}
	copy(uinfo.From[:], ip[:])
	_ = cmsys.StripNoneBig5(uinfo.Nickname[:])

	uinfo.FiveWin = user.FiveWin
	uinfo.FiveLose = user.FiveLose
	uinfo.FiveTie = user.FiveTie
	uinfo.ChcWin = user.ChcWin
	uinfo.ChcLose = user.ChcLose
	uinfo.ChcTie = user.ChcTie
	uinfo.ChessEloRating = user.ChessEloRating
	uinfo.GoWin = user.GoWin
	uinfo.GoLose = user.GoLose
	uinfo.GoTie = user.GoTie
	uinfo.DarkWin = user.DarkWin
	uinfo.DarkLose = user.DarkLose
	uinfo.DarkTie = user.DarkTie
	uinfo.Invisible = user.Invisible && !uinfo.UserLevel.HasUserPerm(ptttype.PERM_VIOLATELAW)
	uinfo.Pager = user.Pager
	uinfo.WithMe = user.WithMe & ^ptttype.WITHME_ALLFLAG
	if (user.WithMe & (user.WithMe << 1) & (ptttype.WITHME_ALLFLAG << 1)) != 0 {
		uinfo.WithMe = 0
	}

	return uinfo
}

func doAloha(utmpID ptttype.UtmpID, uinfo *ptttype.UserInfoRaw, hello []byte) {
	if !ptttype.IS_UTMP {
		return
	}

	filename, err := path.SetHomeFile(&uinfo.UserID, ptttype.FN_ALOHA)
	if err != nil {
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var line []byte
	friendID := &ptttype.UserID_t{}
	for line, err = types.ReadLine(reader); err == nil; line, err = types.ReadLine(reader) {
		//no need to do chomp because it's already taken care in ReadLine
		friendID.CopyFrom(line)
		friendUtmpID, friendInfo := cache.SearchUListUserID(friendID)
		if friendInfo == nil {
			continue
		}

		if types.Cstrcasecmp(friendInfo.UserID[:], uinfo.UserID[:]) == 0 {
			continue
		}

		if !isVisible(friendUtmpID, friendInfo, utmpID, uinfo) {
			continue
		}

		myWrite(utmpID, uinfo, friendInfo.Pid, hello, ptttype.WATERBALL_ALOHA, friendUtmpID, friendInfo)
	}
}
