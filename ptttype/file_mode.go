package ptttype

type FileMode uint8

const (
	FILE_NONE      FileMode = 0x00
	FILE_LOCAL     FileMode = 0x01 /* local saved,  non-mail */
	FILE_READ      FileMode = 0x01 /* already read, mail only */
	FILE_MARKED    FileMode = 0x02 /* non-mail + mail */
	FILE_DIGEST    FileMode = 0x04 /* digest,       non-mail */
	FILE_REPLIED   FileMode = 0x04 /* replied,      mail only */
	FILE_BOTTOM    FileMode = 0x08 /* push_bottom,  non-mail */
	FILE_MULTI     FileMode = 0x08 /* multi send,   mail only */
	FILE_SOLVED    FileMode = 0x10 /* problem solved, sysop/BM non-mail only */
	FILE_HIDE      FileMode = 0x20 /* hide,	in announce */
	FILE_BID       FileMode = 0x20 /* bid,		in non-announce */
	FILE_BM        FileMode = 0x40 /* BM only,	in announce */
	FILE_VOTE      FileMode = 0x40 /* for vote,	in non-announce */
	FILE_ANONYMOUS FileMode = 0x80 /* anonymous file */
)
