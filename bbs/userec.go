package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
type Userec struct {
	Version  uint32  `json:"version"`
	UUserID  UUserID `json:"user_id"`
	Username string  `json:"username"`
	Realname []byte  `json:"realname"`
	Nickname []byte  `json:"nickname"`
	// Pad1       byte

	Uflag ptttype.UFlag `json:"flag"`
	// Unused1      uint32      /* 從前放習慣2, 使用前請先清0 */
	Userlevel    ptttype.PERM `json:"perm"`
	Numlogindays uint32       `json:"login_days"`
	Numposts     uint32       `json:"posts"`
	Firstlogin   types.Time4  `json:"first_login"`
	Lastlogin    types.Time4  `json:"last_login"`
	Lasthost     string       `json:"last_ip"` // last IPv4.
	Money        int32        `json:"money"`
	// Unused2      [4]byte

	Email string `json:"email"`
	// Address []byte
	Justify []byte `json:"justify"`
	// UnusedBirth [3]uint8  /* 生日 月日年 */
	Over18      bool              `json:"over18"`
	PagerUIType uint8             `json:"pager_ui"` /* 呼叫器界面類別 (was: WATER_*) */
	Pager       ptttype.PagerMode `json:"pager"`    /* 呼叫器狀態 */
	Invisible   bool              `json:"invisible"`
	// Unused4     [2]byte
	Exmailbox uint32 `json:"exmail"`

	// r3968 移出 sizeof(chicken_t)=128 bytes
	// Unused5       [4]byte
	Career []byte `json:"career"`
	// UnusedPhone   Phone_t  /* 電話 */
	// Unused6       uint32   /* 從前放轉換前的 numlogins, 使用前請先清0 */
	// Chkpad1       [44]byte
	Role          uint32      `json:"role"`
	LastSeen      types.Time4 `json:"last_seen"`
	TimeSetAngel  types.Time4 `json:"time_set_angel"`
	TimePlayAngel types.Time4 `json:"time_play_angel"`

	LastSong  types.Time4 `json:"last_song"`
	LoginView uint32      `json:"login_view"`
	// Unused8   uint8       // was: channel
	// Pad2 uint8

	Vlcount   uint16 `json:"violation"`
	FiveWin   uint16 `json:"five_win"`
	FiveLose  uint16 `json:"five_lose"`
	FiveTie   uint16 `json:"five_tie"`
	ChcWin    uint16 `json:"chc_win"`
	ChcLose   uint16 `json:"chc_lose"`
	ChcTie    uint16 `json:"chc_tie"`
	Conn6Win  uint16 `json:"conn6_win"`
	Conn6Lose uint16 `json:"conn6_lose"`
	Conn6Tie  uint16 `json:"conn6_tie"`
	// UnusedMind [2]byte /* 舊心情 */
	GoWin     uint16 `json:"go_win"`
	GoLose    uint16 `json:"go_lose"`
	GoTie     uint16 `json:"go_tie"`
	DarkWin   uint16 `json:"dark_win"`
	DarkLose  uint16 `json:"dark_lose"`
	DarkTie   uint16 `json:"dark_tie"` /* 暗棋戰績 和 */
	UaVersion uint8  `json:"ua_version"`

	Signature uint8 `json:"signature"` /* 慣用簽名檔 */
	// Unused10  uint8    /* 從前放好文章數, 使用前請先清0 */
	BadPost uint8   `json:"bad_post"` /* 評價為壞文章數 */
	MyAngel UUserID `json:"angel"`    /* 我的小天使 */
	// Pad3    byte

	ChessEloRating    uint16           `json:"check_rank"`
	WithMe            ptttype.WithMe_t `json:"withme"`
	TimeRemoveBadPost types.Time4      `json:"time_remove_bad_post"`
	TimeViolateLaw    types.Time4      `json:"time_violate_law"`

	UserLevel2 ptttype.PERM2 `json:"user_level2"`
	UpdateTS2  types.Time4   `json:"update_ts2"`
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
