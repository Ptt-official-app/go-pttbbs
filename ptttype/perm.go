package ptttype

type PERM uint32

const (
	PERM_INVALID       PERM = 0o00000000000
	PERM_BASIC         PERM = 0o00000000001 /* 基本權力       */
	PERM_CHAT          PERM = 0o00000000002 /* 進入聊天室     */
	PERM_PAGE          PERM = 0o00000000004 /* 找人聊天       */
	PERM_POST          PERM = 0o00000000010 /* 發表文章       */
	PERM_LOGINOK       PERM = 0o00000000020 /* 註冊程序認證   */
	PERM_MAILLIMIT     PERM = 0o00000000040 /* 信件無上限     */
	PERM_CLOAK         PERM = 0o00000000100 /* 目前隱形中     */
	PERM_SEECLOAK      PERM = 0o00000000200 /* 看見忍者       */
	PERM_XEMPT         PERM = 0o00000000400 /* 永久保留帳號   */
	PERM_SYSOPHIDE     PERM = 0o00000001000 /* 站長隱身術     */
	PERM_BM            PERM = 0o00000002000 /* 板主           */
	PERM_ACCOUNTS      PERM = 0o00000004000 /* 帳號總管       */
	PERM_CHATROOM      PERM = 0o00000010000 /* 聊天室總管     */
	PERM_BOARD         PERM = 0o00000020000 /* 看板總管       */
	PERM_SYSOP         PERM = 0o00000040000 /* 站長           */
	PERM_BBSADM        PERM = 0o00000100000 /* BBSADM         */
	PERM_NOTOP         PERM = 0o00000200000 /* 不列入排行榜   */
	PERM_VIOLATELAW    PERM = 0o00000400000 /* 違法通緝中     */
	PERM_ANGEL         PERM = 0o00001000000 /* 有資格擔任小天使 */
	PERM_NOREGCODE     PERM = 0o00002000000 /* 不允許認證碼註冊 */
	PERM_VIEWSYSOP     PERM = 0o00004000000 /* 視覺站長       */
	PERM_LOGUSER       PERM = 0o00010000000 /* 觀察使用者行蹤 */
	PERM_NOCITIZEN     PERM = 0o00020000000 /* 搋奪公權       */
	PERM_SYSSUPERSUBOP PERM = 0o00040000000 /* 群組長         */
	PERM_ACCTREG       PERM = 0o00100000000 /* 帳號審核組     */
	PERM_PRG           PERM = 0o00200000000 /* 程式組         */
	PERM_ACTION        PERM = 0o00400000000 /* 活動組         */
	PERM_PAINT         PERM = 0o01000000000 /* 美工組         */
	PERM_POLICE_MAN    PERM = 0o02000000000 /* 警察總管       */
	PERM_SYSSUBOP      PERM = 0o04000000000 /* 小組長         */
	PERM_OLDSYSOP      PERM = 0o10000000000 /* 退休站長       */
	PERM_POLICE        PERM = 0o20000000000 /* 警察 */
	// 32 個已經全部用光了。 後面沒有了。
)

const (
	NUMPERMS = 32
)

const (
	PERM_DEFAULT    PERM = PERM_BASIC | PERM_CHAT | PERM_PAGE
	PERM_MANAGER    PERM = PERM_ACCTREG | PERM_ACTION | PERM_PAINT
	PERM_ADMIN      PERM = PERM_ACCOUNTS | PERM_BOARD | PERM_SYSOP | PERM_SYSSUBOP | PERM_SYSSUPERSUBOP | PERM_MANAGER
	PERM_LOGINCLOAK PERM = PERM_SYSOP | PERM_ACCOUNTS
	PERM_SEEULEVELS PERM = PERM_SYSOP
	PERM_SEEBLEVELS PERM = PERM_SYSOP | PERM_BM
	PERM_NOTIMEOUT  PERM = PERM_SYSOP
	PERM_READMAIL   PERM = PERM_BASIC
	PERM_FORWARD    PERM = PERM_LOGINOK /* to do the forwarding */
	PERM_INTERNET   PERM = PERM_LOGINOK /* 身份認證過關的才能寄信到 Internet */
)

func (p PERM) HasUserPerm(perm PERM) bool {
	return p&perm != 0
}

func (p PERM) HasBasicUserPerm(perm PERM) bool {
	return p.HasUserPerm(PERM_BASIC) && p.HasUserPerm(perm)
}

func (p PERM) Hide() bool {
	return p == PERM_SYSOPHIDE
}
