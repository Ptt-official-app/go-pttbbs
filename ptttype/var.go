package ptttype

var (
	STR_REPLY          = []byte("Re:")
	STR_FORWARD        = []byte("Fw:")
	STR_LEGACY_FORWARD = []byte{0x5b, 0xc2, 0xe0, 0xbf, 0xfd, 0x5d}
	STR_SYSOP          = []byte("SYSOP")

	ARTICLE_CLASS_REPLY   = []byte("R:")
	ARTICLE_CLASS_FORWARD = []byte{0xc2, 0xe0}

	STR_DOTS = []byte{0xa1, 0x4b}
)
