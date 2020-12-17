package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
type Userec struct {
	Version  uint32
	UUserID  UUserID
	Userid   string
	Realname string
	Nickname string
	//Pad1       byte

	Uflag ptttype.UFlag
	//Unused1      uint32      /* 從前放習慣2, 使用前請先清0 */
	Userlevel    ptttype.PERM
	Numlogindays uint32
	Numposts     uint32
	Firstlogin   types.Time4
	Lastlogin    types.Time4
	Lasthost     string
	Money        int32
	//Unused2      [4]byte

	Email   string
	Address string
	Justify string
	//UnusedBirth [3]uint8  /* 生日 月日年 */
	Over18      bool
	PagerUIType uint8             /* 呼叫器界面類別 (was: WATER_*) */
	Pager       ptttype.PagerMode /* 呼叫器狀態 */
	Invisible   uint8
	//Unused4     [2]byte
	Exmailbox uint32

	// r3968 移出 sizeof(chicken_t)=128 bytes
	//Unused5       [4]byte
	Career string
	//UnusedPhone   Phone_t  /* 電話 */
	//Unused6       uint32   /* 從前放轉換前的 numlogins, 使用前請先清0 */
	//Chkpad1       [44]byte
	Role          uint32
	LastSeen      types.Time4
	TimeSetAngel  types.Time4
	TimePlayAngel types.Time4

	LastSong  types.Time4
	LoginView uint32
	Unused8   uint8 // was: channel
	Pad2      uint8

	Vlcount   uint16
	FiveWin   uint16
	FiveLose  uint16
	FiveTie   uint16
	ChcWin    uint16
	ChcLose   uint16
	ChcTie    uint16
	Conn6Win  uint16
	Conn6Lose uint16
	Conn6Tie  uint16
	//UnusedMind [2]byte /* 舊心情 */
	GoWin     uint16
	GoLose    uint16
	GoTie     uint16
	DarkWin   uint16
	DarkLose  uint16
	UaVersion uint8

	Signature uint8 /* 慣用簽名檔 */
	//Unused10  uint8    /* 從前放好文章數, 使用前請先清0 */
	BadPost uint8  /* 評價為壞文章數 */
	DarkTie uint16 /* 暗棋戰績 和 */
	MyAngel string /* 我的小天使 */
	//Pad3    byte
}

func NewUserecFromRaw(uid ptttype.Uid, userecRaw *ptttype.UserecRaw) *Userec {
	user := &Userec{}
	user.UUserID = ToUUserID(uid, &userecRaw.UserID)
	user.Version = userecRaw.Version
	user.Userid = types.CstrToString(userecRaw.UserID[:])
	user.Realname = types.Big5ToUtf8(userecRaw.RealName[:])
	user.Nickname = types.Big5ToUtf8(userecRaw.Nickname[:])

	user.Uflag = userecRaw.UFlag
	user.Userlevel = userecRaw.UserLevel
	user.Numlogindays = userecRaw.NumLoginDays
	user.Numposts = userecRaw.NumPosts
	user.Firstlogin = userecRaw.FirstLogin
	user.Lastlogin = userecRaw.LastLogin
	user.Lasthost = types.CstrToString(userecRaw.LastHost[:])
	user.Money = userecRaw.Money
	user.Email = types.CstrToString(userecRaw.Email[:])
	user.Address = types.Big5ToUtf8(userecRaw.Address[:])
	user.Justify = types.Big5ToUtf8(userecRaw.Justify[:])
	user.Over18 = userecRaw.Over18
	user.PagerUIType = userecRaw.PagerUIType
	user.Pager = userecRaw.Pager
	user.Invisible = userecRaw.Invisible
	user.Exmailbox = userecRaw.Exmailbox

	user.Career = types.Big5ToUtf8(userecRaw.Career[:])

	return user
}
