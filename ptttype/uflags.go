package ptttype

type UFlag uint32

const (
	UF_FAV_NOHILIGHT UFlag = 0x00000001 // false if hilight favorite
	UF_FAV_ADDNEW    UFlag = 0x00000002 // true to add new board into one's fav
	// #define UF_PAGER	    0x00000004	// deprecated by cuser.pager: true if pager was OFF last session
	// #define UF_CLOAK	    0x00000008	// deprecated by cuser.invisible: true if cloak was ON last session
	UF_FRIEND         UFlag = 0x00000010 // true if show friends only
	UF_BRDSORT        UFlag = 0x00000020 // true if the boards sorted alphabetical
	UF_ADBANNER       UFlag = 0x00000040 // (was: MOVIE_FLAG, true if show advertisement banner
	UF_ADBANNER_USONG UFlag = 0x00000080 // true if show user songs in banner
	// #define UF_MIND	    0x00000100	// deprecated: true if mind search mode open <-Heat
	UF_DBCS_AWARE       UFlag = 0x00000200 // true if DBCS-aware enabled.
	UF_DBCS_NOINTRESC   UFlag = 0x00000400 // no Escapes interupting DBCS characters
	UF_DBCS_DROP_REPEAT UFlag = 0x00000800 // detect and drop repeated input from evil clients
	// #define UF_DBCS_???	    0x00000800	// reserved
	UF_NO_MODMARK      UFlag = 0x00001000 // true if modified files are NOT marked
	UF_COLORED_MODMARK UFlag = 0x00002000 // true if mod-mark is coloured.
	// #define UF_MODMARK_???   0x00004000	// reserved
	// #define UF_MODMARK_???   0x00008000	// reserved
	UF_DEFBACKUP       UFlag = 0x00010000 // true if user defaults to backup
	UF_NEW_ANGEL_PAGER UFlag = 0x00020000 // true if user (angel) wants the new pager
	UF_REJ_OUTTAMAIL   UFlag = 0x00040000 // true if don't accept outside mails
	UF_SECURE_LOGIN    UFlag = 0x00080000 // true if login from insecure (ex, telnet) connection will be rejected.
	UF_FOREIGN         UFlag = 0x00100000 // true if a foreign
	UF_LIVERIGHT       UFlag = 0x00200000 // true if get "liveright" already
	// #define UF_COUNTRY_???   0x00400000	// reserved
	// #define UF_COUNTRY_???   0x00800000	// reserved
	UF_MENU_LIGHTBAR UFlag = 0x01000000 // true to use lightbar-based menu
	UF_CURSOR_ASCII  UFlag = 0x02000000 // true to enable ASCII-safe cursor.
	// #define UF_???	    0x04000000	// reserved
	// #define UF_???	    0x08000000	// reserved
	// #define UF_???	    0x10000000	// reserved
	// #define UF_???	    0x20000000	// reserved
	// #define UF_???	    0x40000000	// reserved
	// #define UF_???	    0x80000000	// reserved
)
