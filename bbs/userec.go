package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
type Userec struct {
	Version  uint32
	UUserID  UUserID
	Username string
	Realname []byte
	Nickname []byte
	//Pad1       byte

	Uflag ptttype.UFlag
	//Unused1      uint32      /* 從前放習慣2, 使用前請先清0 */
	Userlevel    ptttype.PERM
	Numlogindays uint32
	Numposts     uint32
	Firstlogin   types.Time4
	Lastlogin    types.Time4
	Lasthost     string //last IPv4.
	Money        int32
	//Unused2      [4]byte

	Email string
	//Address []byte
	Justify []byte
	//UnusedBirth [3]uint8  /* 生日 月日年 */
	Over18      bool
	PagerUIType uint8             /* 呼叫器界面類別 (was: WATER_*) */
	Pager       ptttype.PagerMode /* 呼叫器狀態 */
	Invisible   bool
	//Unused4     [2]byte
	Exmailbox uint32

	// r3968 移出 sizeof(chicken_t)=128 bytes
	//Unused5       [4]byte
	Career []byte
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
	DarkTie   uint16 /* 暗棋戰績 和 */
	UaVersion uint8

	Signature uint8 /* 慣用簽名檔 */
	//Unused10  uint8    /* 從前放好文章數, 使用前請先清0 */
	BadPost uint8   /* 評價為壞文章數 */
	MyAngel UUserID /* 我的小天使 */
	//Pad3    byte

	ChessEloRating    uint16
	WithMe            ptttype.WithMe_t
	TimeRemoveBadPost types.Time4
	TimeViolateLaw    types.Time4

	UserLevel2 ptttype.PERM2
	UpdateTS2  types.Time4
}

func NewUserecFromRaw(userecRaw *ptttype.UserecRaw, userec2Raw *ptttype.Userec2Raw) (user *Userec) {
	return &Userec{
		UUserID:  ToUUserID(&userecRaw.UserID),
		Version:  userecRaw.Version,
		Username: types.CstrToString(userecRaw.UserID[:]),
		Realname: types.CstrToBytes(userecRaw.RealName[:]),
		Nickname: types.CstrToBytes(userecRaw.Nickname[:]),

		Uflag:        userecRaw.UFlag,
		Userlevel:    userecRaw.UserLevel,
		Numlogindays: userecRaw.NumLoginDays,
		Numposts:     userecRaw.NumPosts,
		Firstlogin:   userecRaw.FirstLogin,
		Lastlogin:    userecRaw.LastLogin,
		Lasthost:     types.CstrToString(userecRaw.LastHost[:]),
		Money:        userecRaw.Money,
		Email:        types.CstrToString(userecRaw.Email[:]),
		Justify:      types.CstrToBytes(userecRaw.Justify[:]),
		Over18:       userecRaw.Over18,
		PagerUIType:  userecRaw.PagerUIType,
		Pager:        userecRaw.Pager,
		Invisible:    userecRaw.Invisible,
		Exmailbox:    userecRaw.Exmailbox,

		Career: types.CstrToBytes(userecRaw.Career[:]),

		Role:          userecRaw.Role,
		LastSeen:      userecRaw.LastSeen,
		TimeSetAngel:  userecRaw.TimeSetAngel,
		TimePlayAngel: userecRaw.TimePlayAngel,

		LastSong:  userecRaw.LastSong,
		LoginView: userecRaw.LoginView,

		Vlcount:   userecRaw.VlCount,
		FiveWin:   userecRaw.FiveWin,
		FiveLose:  userecRaw.FiveLose,
		FiveTie:   userecRaw.FiveTie,
		ChcWin:    userecRaw.ChcWin,
		ChcLose:   userecRaw.ChcLose,
		ChcTie:    userecRaw.ChcTie,
		Conn6Win:  userecRaw.Conn6Win,
		Conn6Lose: userecRaw.Conn6Lose,
		Conn6Tie:  userecRaw.Conn6Tie,
		GoWin:     userecRaw.GoWin,
		GoLose:    userecRaw.GoLose,
		GoTie:     userecRaw.GoTie,
		DarkWin:   userecRaw.DarkWin,
		DarkLose:  userecRaw.DarkLose,
		DarkTie:   userecRaw.DarkTie,

		UaVersion: userecRaw.UaVersion,

		Signature: userecRaw.Signature,
		BadPost:   userecRaw.BadPost,
		MyAngel:   ToUUserID(&userecRaw.MyAngel),

		ChessEloRating:    userecRaw.ChessEloRating,
		WithMe:            userecRaw.WithMe,
		TimeRemoveBadPost: userecRaw.TimeRemoveBadPost,
		TimeViolateLaw:    userecRaw.TimeViolateLaw,

		UserLevel2: userec2Raw.UserLevel2,
		UpdateTS2:  userec2Raw.UpdateTS,
	}
}
