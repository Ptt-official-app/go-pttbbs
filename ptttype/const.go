package ptttype

const (
	//////////
	//pttstruct.h
	//
	//SHMSIZE is computed in cache/shm.go NewSHM
	//////////
	IDLEN   = 12 /* Length of bid/uid */
	IPV4LEN = 15 /* a.b.c.d form */

	PASS_INPUT_LEN = 8 /* Length of valid input password length.
	   For DES, set to 8. */
	PASSLEN = 14 /* Length of encrypted passwd field */
	REGLEN  = 38 /* Length of registration data */

	REALNAMESZ = 20 /* Size of real-name field */
	NICKNAMESZ = 24 /* SIze of nick-name field */
	EMAILSZ    = 50 /* Size of email field */
	ADDRESSSZ  = 50 /* Size of address field */
	CAREERSZ   = 40 /* Size of career field */
	PHONESZ    = 20 /* Size of phone field */

	USERNAMESZ = 24 /* Size of Username in MailQueue */
	RCPTSZ     = 50 /* Size of RCPT in MailQueue */

	PASSWD_VERSION = 4194

	BTLEN = 48 /* Length of board title */

	TTLEN = 64 /* Length of title */
	FNLEN = 28 /* Length of filename */

	// Prefix of str_reply and str_forward usually needs longer.
	// In artitlcle list, the normal length is (80-34)=46.
	DISP_TTLEN = 46

	STRLEN = 80 /* Length of most string data */

	FAVMAX   = 1024 /* Max boards of Myfavorite */
	FAVGMAX  = 32   /* Max groups of Myfavorite */
	FAVGSLEN = 8    /* Max Length of Description String */

	ONEKEY_SIZE = int('~')

	ANSILINELEN = 511 /* Maximum Screen width in chars */

	/* USHM_SIZE 比 MAX_ACTIVE 大是為了防止檢查人數上限時, 又同時衝進來
	 * 會造成找 shm 空位的無窮迴圈.
	 * 又, 因 USHM 中用 hash, 空間稍大時效率較好. */
	USHM_SIZE = MAX_ACTIVE * 41 / 40

	MAX_BMs = 4
)

const (
	//mbbsd/register.c line: 381
	FN_FRESH = ".fresh"

	//mbbsd/register.c line: 415
	CLEAN_USER_EXPIRE_RANGE_MIN = 365 * 12 * 60 // 180 days.

	//mbbsd/user.c line: 42
	DIR_TMP  = "tmp"
	DIR_HOME = "home"
)
