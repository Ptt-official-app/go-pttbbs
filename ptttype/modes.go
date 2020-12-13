package ptttype

type PagerMode uint8

const (
	PAGER_OFF PagerMode = iota
	PAGER_ON
	PAGER_DISABLE
	PAGER_ANTIWB
	PAGER_FRIENDONLY
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
