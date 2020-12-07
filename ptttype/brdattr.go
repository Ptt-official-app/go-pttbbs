package ptttype

import "unsafe"

type BrdAttr uint32

const (
	// TODO BRD 快爆了，怎麼辦？ 準備從 pad3 偷一個來當 attr2 吧...
	// #define BRD_NOZAP       	0x00000001	/* 不可 ZAP */
	BRD_NOCOUNT BrdAttr = 0x00000002 /* 不列入統計 */
	//#define BRD_NOTRAN		0x00000004	/* 不轉信 */
	BRD_GROUPBOARD       BrdAttr = 0x00000008 /* 群組板 */
	BRD_HIDE             BrdAttr = 0x00000010 /* 隱藏板 (看板好友才可看) */
	BRD_POSTMASK         BrdAttr = 0x00000020 /* 限制發表或閱讀 */
	BRD_ANONYMOUS        BrdAttr = 0x00000040 /* 匿名板 */
	BRD_DEFAULTANONYMOUS BrdAttr = 0x00000080 /* 預設匿名板 */
	BRD_NOCREDIT         BrdAttr = 0x00000100 /* 發文無獎勵看板 */
	BRD_VOTEBOARD        BrdAttr = 0x00000200 /* 連署機看板 */
	BRD_WARNEL           BrdAttr = 0x00000400 /* 連署機看板 */
	BRD_TOP              BrdAttr = 0x00000800 /* 熱門看板群組 */
	BRD_NORECOMMEND      BrdAttr = 0x00001000 /* 不可推薦 */
	BRD_ANGELANONYMOUS   BrdAttr = 0x00002000 /* 小天使可匿名 */
	BRD_BMCOUNT          BrdAttr = 0x00004000 /* 板主設定列入記錄 */
	BRD_SYMBOLIC         BrdAttr = 0x00008000 /* symbolic link to board */
	BRD_NOBOO            BrdAttr = 0x00010000 /* 不可噓 */
	//BRD_LOCALSAVE		0x00020000	/* 預設 Local Save */
	BRD_RESTRICTEDPOST  BrdAttr = 0x00040000 /* 板友才能發文 */
	BRD_GUESTPOST       BrdAttr = 0x00080000 /* guest能 post */
	BRD_COOLDOWN        BrdAttr = 0x00100000 /* 冷靜 */
	BRD_CPLOG           BrdAttr = 0x00200000 /* 自動留轉錄記錄 */
	BRD_NOFASTRECMD     BrdAttr = 0x00400000 /* 禁止快速推文 */
	BRD_IPLOGRECMD      BrdAttr = 0x00800000 /* 推文記錄 IP */
	BRD_OVER18          BrdAttr = 0x01000000 /* 十八禁 */
	BRD_NOREPLY         BrdAttr = 0x02000000 /* 不可回文 */
	BRD_ALIGNEDCMT      BrdAttr = 0x04000000 /* 對齊式的推文 */
	BRD_NOSELFDELPOST   BrdAttr = 0x08000000 /* 不可自刪 */
	BRD_BM_MASK_CONTENT BrdAttr = 0x10000000 /* 允許板主刪除特定文字 */
)

const BRD_ATTR_SZ = unsafe.Sizeof(BrdAttr(0))
