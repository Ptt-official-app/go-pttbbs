package ptttype

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

type BoardHeaderRaw struct {
	//Require updating SHM_VERSION if BOARD_HEADER_RAW_SZ is changed.
	Brdname            BoardID_t /* bid */
	Title              BoardTitle_t
	BM                 BM_t /* BMs' userid, token '/' */
	Pad1               [3]byte
	BrdAttr            uint32            /* board的屬性 */
	ChessCountry       byte              /* 棋國 */
	VoteLimitPosts_    uint8             /* (已停用) 連署 : 文章篇數下限 */
	VoteLimitLogins    uint8             /* 連署 : 登入次數下限 */
	Pad2_1             [1]uint8          /* (已停用) 連署 : 註冊時間限制 */
	BUpdate            types.Time4       /* note update time */
	PostLimitPosts_    uint8             /* (已停用) 發表文章 : 文章篇數下限 */
	PostLimitLogins    uint8             /* 發表文章 : 登入次數下限 */
	Pad2_2             [1]uint8          /* (已停用) 發表文章 : 註冊時間限制 */
	BVote              uint8             /* 正舉辦 Vote 數 */
	VTime              types.Time4       /* Vote close time */
	Level              uint32            /* 可以看此板的權限 */
	PermReload         types.Time4       /* 最後設定看板的時間 */
	Gid                Gid               /* 看板所屬的類別 ID */
	Next               [2]Bid            /* 在同一個gid下一個看板 動態產生*/
	FirstChild         [BSORT_BY_MAX]Bid /* 屬於這個看板的第一個子看板 */
	Parent             Bid               /* 這個看板的 parent 看板 bid */
	ChildCount         int32             /* 有多少個child */
	NUser              int32             /* 多少人在這板 */
	PostExpire         int32             /* postexpire */
	EndGamble          types.Time4
	PostType           [33]byte
	PostTypeF          byte
	FastRecommendPause uint8 /* 快速連推間隔 */
	VoteLimitBadpost   uint8 /* 連署 : 劣文上限 */
	PostLimitBadpost   uint8 /* 發表文章 : 劣文上限 */
	Pad3               [3]byte
	SRexpire           types.Time4 /* SR Records expire time */
	Pad4               [40]byte
}

//Require updating SHM_VERSION if BOARD_HEADER_RAW_SZ is changed.

var EMPTY_BOARD_HEADER_RAW = BoardHeaderRaw{}

const BOARD_HEADER_RAW_SZ = unsafe.Sizeof(EMPTY_BOARD_HEADER_RAW)
const BOARD_HEADER_BRDNAME_OFFSET = unsafe.Offsetof(EMPTY_BOARD_HEADER_RAW.Brdname)
const BOARD_HEADER_TITLE_OFFSET = unsafe.Offsetof(EMPTY_BOARD_HEADER_RAW.Title)

const BOARD_HEADER_FIRST_CHILD_OFFSET = unsafe.Offsetof(EMPTY_BOARD_HEADER_RAW.FirstChild)
