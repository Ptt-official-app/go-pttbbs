package ptttype

import "unsafe"

type PagerMode uint8

const (
	PAGER_OFF PagerMode = iota
	PAGER_ON
	PAGER_DISABLE
	PAGER_ANTIWB
	PAGER_FRIENDONLY

	PAGER_MODES
)

type StateMode uint16

const (
	MODE_STARTED     StateMode = 0x0001 /* 是否已經進入系統 */
	MODE_POST        StateMode = 0x0002 /* 是否可以在 currboard 發表文章 */
	MODE_POSTCHECKED StateMode = 0x0004 /* 是否已檢查在 currboard 發表文章的權限 */
	MODE_BOARD       StateMode = 0x0008 /* 是否可以在 currboard 刪除、mark文章 */
	MODE_GROUPOP     StateMode = 0x0010 /* 是否為小組長 (可以在 MENU 開板) */
	MODE_DIGEST      StateMode = 0x0020 /* 是否為 digest mode */
	MODE_SELECT      StateMode = 0x0080 /* 搜尋使用者標題 */
	MODE_DIRTY       StateMode = 0x0100 /* 是否更動過 userflag */
)

type UserOpMode uint8

const (
	// it's better that we explicitly specify
	// the number in the code,
	// to compare with c-pttbbs.
	USER_OP_IDLE     UserOpMode = 0
	USER_OP_MMENU    UserOpMode = 1 /* menu mode */
	USER_OP_ADMIN    UserOpMode = 2
	USER_OP_MAIL     UserOpMode = 3
	USER_OP_TMENU    UserOpMode = 4
	USER_OP_UMENU    UserOpMode = 5
	USER_OP_XMENU    UserOpMode = 6
	USER_OP_CLASS    UserOpMode = 7
	USER_OP_PMENU    UserOpMode = 8
	USER_OP_NMENU    UserOpMode = 9
	USER_OP_PSALE    UserOpMode = 10
	USER_OP_POSTING  UserOpMode = 11 /* boards & class */
	USER_OP_READBRD  UserOpMode = 12
	USER_OP_READING  UserOpMode = 13
	USER_OP_READNEW  UserOpMode = 14
	USER_OP_SELECT   UserOpMode = 15
	USER_OP_RMAIL    UserOpMode = 16 /* mail menu */
	USER_OP_SMAIL    UserOpMode = 17
	USER_OP_CHATING  UserOpMode = 18 /* talk menu */
	USER_OP_XMODE    UserOpMode = 19
	USER_OP_FRIEND   UserOpMode = 20
	USER_OP_LAUSERS  UserOpMode = 21
	USER_OP_LUSERS   UserOpMode = 22
	USER_OP_MONITOR  UserOpMode = 23
	USER_OP_PAGE     UserOpMode = 24
	USER_OP_TQUERY   UserOpMode = 25
	USER_OP_TALK     UserOpMode = 26
	USER_OP_EDITPLAN UserOpMode = 27 /* user menu */
	USER_OP_EDITSIG  UserOpMode = 28
	USER_OP_VOTING   UserOpMode = 29
	USER_OP_XINFO    UserOpMode = 30
	USER_OP_MSYSOP   UserOpMode = 31
	// USER_OP_WWW             32
	// USER_OP_BIG2            33
	USER_OP_REPLY    UserOpMode = 34
	USER_OP_HIT      UserOpMode = 35
	USER_OP_DBACK    UserOpMode = 36
	USER_OP_NOTE     UserOpMode = 37
	USER_OP_EDITING  UserOpMode = 38
	USER_OP_MAILALL  UserOpMode = 39
	USER_OP_MJ       UserOpMode = 40
	USER_OP_P_FRIEND UserOpMode = 41
	USER_OP_LOGIN    UserOpMode = 42 /* main menu */
	USER_OP_DICT     UserOpMode = 43
	// USER_OP_BRIDGE          44
	USER_OP_ARCHIE  UserOpMode = 45
	USER_OP_GOPHER  UserOpMode = 46
	USER_OP_NEWS    UserOpMode = 47
	USER_OP_LOVE    UserOpMode = 48
	USER_OP_EDITEXP UserOpMode = 49
	USER_OP_IPREG   UserOpMode = 50
	USER_OP_NADM    UserOpMode = 51
	USER_OP_DRINK   UserOpMode = 52
	USER_OP_CAL     UserOpMode = 53
	// USER_OP_PROVERB         54
	USER_OP_ANNOUNCE UserOpMode = 55 /* announce */
	USER_OP_EDNOTE   UserOpMode = 56
	USER_OP_CDICT    UserOpMode = 57
	USER_OP_LOBJ     UserOpMode = 58
	USER_OP_OSONG    UserOpMode = 59
	USER_OP_CHICKEN  UserOpMode = 60
	USER_OP_TICKET   UserOpMode = 61
	// USER_OP_GUESSNUM        62
	USER_OP_AMUSE   UserOpMode = 63
	USER_OP_OTHELLO UserOpMode = 64
	USER_OP_DICE    UserOpMode = 65
	// USER_OP_VICE            66
	// USER_OP_BBCALL          67
	USER_OP_VIOLATELAW UserOpMode = 68
	USER_OP_M_FIVE     UserOpMode = 69
	USER_OP_M_CONN6    UserOpMode = 70
	// USER_OP_TENHALF         71
	// USER_OP_CARD_99         72
	// USER_OP_RAIL_WAY        73
	USER_OP_SREG UserOpMode = 74
	USER_OP_CHC  UserOpMode = 75 /* Chinese chess */
	USER_OP_DARK UserOpMode = 76 /* 中國暗棋 */
	// USER_OP_TMPJACK         77
	// USER_OP_JCEE		78
	USER_OP_REEDIT UserOpMode = 79
	// USER_OP_BLOGGING        80	/* 已停用 */
	USER_OP_CHESSWATCHING    UserOpMode = 81
	USER_OP_UMODE_GO         UserOpMode = 82
	USER_OP_DEBUGSLEEPING    UserOpMode = 83
	USER_OP_UMODE_CONN6      UserOpMode = 84
	USER_OP_REVERSI          UserOpMode = 85
	USER_OP_UMODE_BBSLUA     UserOpMode = 86
	USER_OP_UMODE_ASCIIMOVIE UserOpMode = 87

	USER_OP_MODE_MAX UserOpMode = 88 /* 所有其他選單動態須在此之前 */
)

const USER_OP_MODE_SZ = unsafe.Sizeof(UserOpMode(0))

type WaterBall uint8

const (
	WATERBALL_GENERAL WaterBall = 0
	WATERBALL_PREEDIT WaterBall = 1
	WATERBALL_ALOHA   WaterBall = 2
	WATERBALL_SYSOP   WaterBall = 3
	WATERBALL_CONFIRM WaterBall = 4

	WATERBALL_ANGEL          WaterBall = 5
	WATERBALL_ANSWER         WaterBall = 6
	WATERBALL_CONFIRM_ANGEL  WaterBall = 7
	WATERBALL_CONFIRM_ANSWER WaterBall = 8
)
