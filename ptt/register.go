package ptt

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//Login
//
//adopted from the original start_client with LOGINASNEW.
//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L1399
//https://github.com/ptt/pttbbs/blob/master/mbbsd/mbbsd.c#L714
func Register(
	userID *ptttype.UserID_t,
	passwd []byte,
	fromHost *ptttype.IPv4_t,
	email *ptttype.Email_t,
	isEmailVerified bool,
	isAdbannerUSong bool,

	nickname *ptttype.Nickname_t,
	realname *ptttype.RealName_t,
	career *ptttype.Career_t,
	address *ptttype.Address_t,
	over18 bool,
) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {
	uid, user, err = NewRegister(
		userID,
		passwd,
		fromHost,
		email,
		isEmailVerified,
		isAdbannerUSong,

		nickname,
		realname,
		career,
		address,
		over18,
	)
	if err != nil {
		return 0, nil, err
	}

	err = mkUserDir(&user.UserID)
	if err != nil {
		return 0, nil, err
	}

	err = reginitFav(uid, user)
	if err != nil {
		return 0, nil, err
	}

	err = userLogin(uid, user, fromHost)
	if err != nil {
		return 0, nil, err
	}

	return uid, user, nil
}

//NewRegister
//
//XXX Assuming valid input. need to verify email at api.
func NewRegister(
	userID *ptttype.UserID_t,
	passwd []byte,
	fromHost *ptttype.IPv4_t,
	email *ptttype.Email_t,
	isEmailVerified bool,
	isAdbannerUSong bool,

	nickname *ptttype.Nickname_t,
	realname *ptttype.RealName_t,
	career *ptttype.Career_t,
	address *ptttype.Address_t,
	over18 bool,
) (uid ptttype.Uid, user *ptttype.UserecRaw, err error) {

	//https://github.com/ptt/pttbbs/blob/master/mbbsd/register.c#L723
	if isBadUserID(userID) {
		return 0, nil, ptttype.ErrInvalidUserID
	}

	if isReservedUserID(userID) {
		return 0, nil, ptttype.ErrInvalidUserID
	}

	logrus.Infof("NewRegister: good userID: %v", string(userID[:]))

	// line: 758
	newUser := &ptttype.UserecRaw{}
	newUser.Version = ptttype.PASSWD_VERSION
	newUser.UserLevel = ptttype.PERM_DEFAULT
	newUser.UFlag = ptttype.UF_BRDSORT | ptttype.UF_ADBANNER | ptttype.UF_CURSOR_ASCII
	newUser.FirstLogin = types.NowTS()
	newUser.LastLogin = newUser.FirstLogin
	newUser.Pager = ptttype.PAGER_ON
	newUser.NumLoginDays = 1
	newUser.UaVersion = getSystemUaVersion()
	copy(newUser.LastHost[:], fromHost[:])

	if ptttype.DBCSAWARE {
		newUser.UFlag |= ptttype.UF_DBCS_AWARE | ptttype.UF_DBCS_DROP_REPEAT
	}

	// XXX REQUIRE_VERIFY_EMAIL_AT_REGISTER
	if email != nil {
		copy(newUser.Email[:], email[:])
	}

	// XXX UF_ADBANNER_USONG
	if isAdbannerUSong {
		newUser.UFlag |= ptttype.UF_ADBANNER_USONG
	}

	// line: 857
	passwdHash, err := cmbbs.GenPasswd(passwd)
	if err != nil {
		return 0, nil, err
	}
	copy(newUser.PasswdHash[:], passwdHash[:])

	copy(newUser.Nickname[:], nickname[:])
	copy(newUser.RealName[:], realname[:])
	copy(newUser.Career[:], career[:])
	copy(newUser.Address[:], address[:])
	newUser.Over18 = over18

	if ptttype.REQUIRE_SECURE_CONN_TO_REGISTER {
		newUser.UFlag |= ptttype.UF_SECURE_LOGIN
	}

	copy(newUser.UserID[:], userID[:])

	err = SetupNewUser(newUser)
	if err != nil {
		log.Errorf("register.NewRegister: unable to SetupNewUser: userID: %v e: %v", userID, err)
		return 0, nil, err
	}

	// read and ensure that every works as expected.
	uid, _, err = InitCurrentUser(userID)
	if err != nil {
		log.Errorf("register.NewRegister: unable to initCurrentUser: userID: %v e: %v", userID, err)
		return 0, nil, err
	}

	err = ensureErasingOldUser(uid, userID)
	if err != nil {
		return 0, nil, err
	}

	// XXX no define of USE_REMOVEBM_ON_NEWREG (line: 985)

	// if email verified.
	if isEmailVerified {
		emailErr := registerCheckAndUpdateEmaildb(newUser, &newUser.Email)
		if emailErr == nil {
			justify := ptttype.Reg_t{}
			copy(justify[:ptttype.REGLEN], []byte(fmt.Sprintf("<E-Mail>: %v", types.NowTS().Cdate())))
			err = pwcuRegCompleteJustify(uid, userID, &justify)
			if err != nil {
				return 0, nil, err
			}
		}
	}

	user, err = passwdSyncQuery(uid)
	if err != nil {
		return 0, nil, err
	}
	return uid, user, nil
}

func ensureErasingOldUser(uid ptttype.Uid, userID *ptttype.UserID_t) (err error) {
	filename := path.SetHomePath(userID)
	tmpFilename := filename + fmt.Sprintf(".%v", types.NowTS())
	if !types.IsDir(filename) {
		return nil
	}

	err = os.Rename(filename, tmpFilename)
	if err == nil {
		return nil
	}

	pwcuBitDisableLevel(uid, userID, ptttype.PERM_BASIC)

	return FatalLockedUser(userID)
}

func registerCheckAndUpdateEmaildb(user *ptttype.UserecRaw, email *ptttype.Email_t) (err error) {

	_, err = registerCountEmail(user, email)
	if err != nil {
		return err
	}

	if ptttype.USE_EMAILDB {
		err = emailDBUpdateEmail(&user.UserID, email)
		if err != nil {
			return err
		}
	}

	if ptttype.USE_VERIFYDB {
		lcemail := types.CstrTolower(email[:])
		err = verifyDBSet(&user.UserID, int64(user.FirstLogin), ptttype.VMETHOD_EMAIL, lcemail, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func registerCountEmail(user *ptttype.UserecRaw, email *ptttype.Email_t) (count int, err error) {

	if ptttype.USE_EMAILDB {
		r, err := emailDBCheckEmail(&user.UserID, email)
		if err != nil {
			return r, err
		}
		count = r
	}

	if ptttype.USE_VERIFYDB {
		lcemail := types.CstrTolower(email[:])
		_, countOther, err := verifyDBCountByVerify(ptttype.VMETHOD_EMAIL, lcemail)
		if err != nil {
			return -1, err
		}

		if count < countOther {
			count = countOther
		}
	}

	return count, nil
}

func getSystemUaVersion() uint8 {
	file, err := os.Open(ptttype.HAVE_USERAGREEMENT_VERSION)
	if err != nil {
		return 0
	}
	defer file.Close()

	var version uint
	_, err = fmt.Fscanf(file, "%v", &version)
	if err != nil {
		log.Errorf("getSystemUaVersion: unable to get version: e: %v", err)
		return 0
	}

	if version > 255 {
		version = 255
	}

	return uint8(version)
}

func SetupNewUser(user *ptttype.UserecRaw) error {
	// XXX race-condition, need to setup RWLock (across-processes)
	//
	uid, err := cache.DoSearchUserRaw(&user.UserID, nil)
	if err != nil {
		return err
	}

	if uid != 0 {
		return ptttype.ErrUserIDAlreadyExists
	}

	/* Lazy method : 先找尋已經清除的過期帳號 */
	uid, err = cache.DoSearchUserRaw(&ptttype.EMPTY_USER_ID, nil)
	if err != nil {
		return err
	}

	if uid == 0 { //unable to find empty user.
		err := tryCleanUser()
		if err != nil {
			return err
		}
	}

	//init passwd-semaphores
	//XXX move to init-config as 1-time setup.
	//err = cmbbs.PasswdInit()

	err = cmbbs.PasswdLock()
	if err != nil {
		return err
	}
	defer cmbbs.PasswdUnlock()

	uid, err = cache.DoSearchUserRaw(&ptttype.EMPTY_USER_ID, nil)
	if err != nil {
		return err
	}

	err = cache.SetUserID(uid, &user.UserID)
	if err != nil {
		return err
	}

	_, _ = cache.SetUMoney(uid, user.Money)

	err = passwdSyncUpdate(uid, user)
	if err != nil {
		return err
	}

	return nil
}

func tryCleanUser() error {

	isToCleanUser, err := isToCleanUser()
	if err != nil {
		return err
	}
	if !isToCleanUser {
		return nil
	}

	err = touchFresh()
	if err != nil {
		return err
	}

	/* 不曉得為什麼要從 2 開始... Ptt:因為SYSOP在1 */
	for uid := ptttype.Uid(2); uid <= ptttype.MAX_USERS; uid++ {
		// XXX ignoring err, do the log
		user, err := passwdSyncQuery(uid)
		if err != nil {
			log.Errorf("register.tryCleanUser: unable to PasswdSyncQuery: uid: %v e: %v", uid, err)
		}

		_, err = checkAndExpireAccount(uid, user, ptttype.CLEAN_USER_EXPIRE_RANGE_MIN)
		if err != nil {
			log.Errorf("register.tryCleanUser: unable to checkAndExpireAccount: uid: %v e: %v", uid, err)
		}

	}

	return nil
}

//isToCleanUser
func isToCleanUser() (bool, error) {
	theStat, err := os.Stat(ptttype.FN_FRESH)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	return theStat.ModTime().Before(time.Now().Add(-3600 * types.TS_TO_NANO_TS)), nil
}

//touchFresh
func touchFresh() error {
	file, err := os.OpenFile(ptttype.FN_FRESH, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(time.Now().String())

	return nil
}

func checkAndExpireAccount(uid ptttype.Uid, user *ptttype.UserecRaw, expireRange int) (int, error) {
	expireValue := computeUserExpireValue(user)
	if expireValue >= 0 { // not expired yet.
		return expireValue, nil
	}

	if -expireValue > expireRange {
		killUser(uid, &user.UserID)
	} else {
		expireValue = 0
	}

	return expireValue, nil
}

func computeUserExpireValue(user *ptttype.UserecRaw) int {
	if (user.UserID[0] == byte(0)) ||
		(user.UserLevel&ptttype.PERM_XEMPT) != 0 ||
		types.Cstrcmp(user.UserID[:], ptttype.USER_ID_GUEST[:]) == 0 { // no expire
		return 999999
	}

	valMinute := int(types.NowTS()-user.LastLogin) / 60 // min-time since last login

	if types.Cstrcmp(user.UserID[:], ptttype.USER_ID_REGNEW[:]) == 0 { // allow only 30 minutes when doing new-user.
		return 30 - valMinute
	}

	if user.UserLevel&(ptttype.PERM_LOGINOK|ptttype.PERM_VIOLATELAW) != 0 {
		return ptttype.KEEP_DAYS_REGGED*24*60 - valMinute
	}

	return ptttype.KEEP_DAYS_UNREGGED*24*60 - valMinute
}

func CheckEmailAllowRejectLists(email string) (err error) {
	err = ensureNewestAllowRejectLists()
	if err != nil {
		return err
	}

	//allow
	isAllow := false
	for _, each := range ptttype.ALLOW_EMAIL_LIST {
		isValid, err := each.IsValid(email)
		if err != nil {
			return err
		}
		if isValid {
			isAllow = true
			break
		}
	}
	if !isAllow {
		return ErrInvalidEmail
	}

	for _, each := range ptttype.REJECT_EMAIL_LIST {
		isValid, err := each.IsValid(email)
		if err != nil {
			return err
		}

		if isValid {
			return ErrInvalidEmail
		}
	}

	return nil
}

func ensureNewestAllowRejectLists() (err error) {
	ptttype.ALLOW_EMAIL_LIST, ptttype.ALLOW_EMAIL_LIST_UPDATE_TS, err = ensureNewestAllowRejectListCore(ptttype.FN_ALLOW_EMAIL_LIST, ptttype.ALLOW_EMAIL_LIST, ptttype.ALLOW_EMAIL_LIST_UPDATE_TS)
	if err != nil {
		return err
	}

	ptttype.REJECT_EMAIL_LIST, ptttype.REJECT_EMAIL_LIST_UPDATE_TS, err = ensureNewestAllowRejectListCore(ptttype.FN_REJECT_EMAIL_LIST, ptttype.REJECT_EMAIL_LIST, ptttype.REJECT_EMAIL_LIST_UPDATE_TS)
	if err != nil {
		return err
	}

	return nil
}

func ensureNewestAllowRejectListCore(filename string, origList []*ptttype.AllowRejectEmail, origUpdateTS types.Time4) (newList []*ptttype.AllowRejectEmail, newUpdateNanoTS types.Time4, err error) {
	theStat, err := os.Stat(filename)

	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, 0, err
	}

	mTime := types.TimeToTime4(theStat.ModTime())
	if mTime <= origUpdateTS {
		return origList, origUpdateTS, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for line, err := types.ReadLine(reader); err == nil; line, err = types.ReadLine(reader) {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		each := ptttype.NewAllowRejectEmail(string(line))

		if each == nil {
			continue
		}
		newList = append(newList, each)
	}
	return newList, mTime, nil
}

func isBadUserID(userID *ptttype.UserID_t) bool {
	if !userID.IsValid() {
		return true
	}
	if types.Cstrcmp(userID[:], []byte(ptttype.STR_REGNEW)) == 0 {
		return true
	}
	if types.Cstrcmp(userID[:], []byte(ptttype.STR_GUEST)) == 0 {
		return true
	}

	return false
}

func isReservedUserID(userID *ptttype.UserID_t) bool {
	for _, each := range ptttype.ReservedUserIDs {
		logrus.Infof("to compare: each: %v userID: %v", string(each), string(userID[:]))
		if types.Cstrcasecmp(userID[:], each[:]) == 0 {
			return true
		}
	}

	return false
}
