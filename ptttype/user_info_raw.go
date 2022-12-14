package ptttype

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

type UserInfoRaw struct {
	// Require updating SHM_VERSION if USER_INFO_RAW_SZ is changed.
	UID      UID         /* Used to find user name in passwd file XXX need to check whether it is Uid or UidInStore */
	Pid      types.Pid_t /* kill() to notify user of talk request */
	SockAddr int32       /* ... */

	/* user data */
	UserLevel  PERM
	UserID     UserID_t
	Nickname   Nickname_t
	From       From_t         /* machine name the user called in from */
	FromIP     types.InAddr_t // was: int     from_alias;
	DarkWin    uint16
	DarkLose   uint16
	Gap0       byte
	AngelPause uint8 // TODO move to somewhere else in future.
	DarkTie    uint16

	/* friends */
	FriendTotal int32 /* 好友比較的cache 大小 */
	NFriends    int16 /* 下面 friend[] 只用到前幾個,
	   用來 bsearch */
	Unused3_     int16
	MyFriend     [MAX_FRIEND]int32
	Gap1         [4]byte
	FriendOnline [MAX_FRIEND]FriendOnline /* point到線上好友 utmpshm的位置 */
	/* 好友比較的cache 前兩個bit是狀態 */
	Gap2   [4]byte
	Reject [MAX_REJECT]int32
	Gap3   [4]byte

	/* messages */
	MsgCount uint8
	Unused4_ [3]byte
	Msgs     [MAX_MSGS]MsgQueueRaw
	Gap4     [MSG_QUEUE_RAW_SZ] /* avoid msgs racing and overflow */ byte

	/* user status */
	Birth        int8       /* 是否是生日 Ptt*/
	Active       uint8      /* When allocated this field is true */
	Invisible    bool       /* Used by cloaking function in Xyz menu */
	Mode         UserOpMode /* UL/DL, Talk Mode, Chat Mode, ... */
	Pager        PagerMode  /* pager toggle, YEA, or NA */
	Unused5_     byte
	Conn6Win     uint16
	LastAct      types.Time4 /* 上次使用者動的時間 */
	Alerts       byte        /* mail alert, passwd update... */
	UnusedMind_  byte
	Conn6Lose    uint16
	UnusedMind2_ byte

	/* chatroom/talk/games calling */
	Sig        byte /* signal type */
	Conn6Tie   uint16
	DestUID    int32 /* talk uses this to identify who called */
	DestUip    int32 /* dest index in utmpshm->uinfo[] */
	SockActive uint8 /* Used to coordinate talk requests */

	/* chat */
	InChat uint8    /* for in_chat commands   */
	Chatid ChatID_t /* chat id, if in chat mode */

	/* games */
	LockMode uint8    /* 不准 multi_login 玩的東西 */
	Turn     byte     /* 遊戲的先後 */
	Mateid   UserID_t /* 遊戲對手的 id */
	Color    byte     /* 暗棋 顏色 */

	/* game record */
	FiveWin        uint16
	FiveLose       uint16
	FiveTie        uint16
	ChcWin         uint16
	ChcLose        uint16
	ChcTie         uint16
	ChessEloRating uint16
	GoWin          uint16
	GoLose         uint16
	GoTie          uint16

	/* misc */
	WithMe WithMe_t
	BrcID  uint32

	// XXX affected by NOKILLWATERBALL
	WBTime types.Time4
}

// Require updating SHM_VERSION if USER_INFO_RAW_SZ is changed.
var EMPTY_USER_INFO_RAW = UserInfoRaw{}

const (
	USER_INFO_RAW_SZ         = unsafe.Sizeof(EMPTY_USER_INFO_RAW)
	USER_INFO_USER_ID_OFFSET = unsafe.Offsetof(EMPTY_USER_INFO_RAW.UserID)
	USER_INFO_PID_OFFSET     = unsafe.Offsetof(EMPTY_USER_INFO_RAW.Pid)
	USER_INFO_MODE_OFFSET    = unsafe.Offsetof(EMPTY_USER_INFO_RAW.Mode)
	USER_INFO_MODE_SZ        = unsafe.Sizeof(EMPTY_USER_INFO_RAW.Mode)
)

type FriendOnline uint32

func (f FriendOnline) ToUtmpID() UtmpID {
	return UtmpID(f & 0xffffff)
}

func (f FriendOnline) ToFriendStat() FriendStat {
	return FriendStat(f >> 24)
}
