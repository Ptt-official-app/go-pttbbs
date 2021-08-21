package ptt

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

const (
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/acl.c
	BAKUMAN_OBJECT_TYPE_USER  = "u"
	BAKUMAN_OBJECT_TYPE_BOARD = "b"
)

var (
	BAKUMAN_REASON_LEN = ptttype.BTLEN

	CROSS_POST_HIDDEN_BOARD  = []byte("\xa1\xb0 [\xa5\xbb\xa4\xe5\xc2\xe0\xbf\xfd\xa6\xdb\xacY\xc1\xf4\xa7\xce\xac\xdd\xaaO]\n\n")
	CROSS_POST_BOARD_PREFIX  = []byte("\xa1\xb0 [\xa5\xbb\xa4\xe5\xc2\xe0\xbf\xfd\xa6\xdb ")
	CROSS_POST_BOARD_INFIX   = []byte(" \xac\xdd\xaaO #")
	CROSS_POST_BOARD_POSTFIX = []byte("]\n\n")

	CROSS_POST_COMMENT_PREFIX       = "\xa1\xb0 \x1b[1;32m"
	CROSS_POST_COMMENT_INFIX        = "\x1b[0;32m:\xc2\xe0\xbf\xfd\xa6\xdc"
	CROSS_POST_COMMENT_HIDDEN_BOARD = "\xacY\xc1\xf4\xa7\xce\xac\xdd\xaaO"
	CROSS_POST_COMMENT_BOARD        = "\xac\xdd\xaaO"

	INVIS_BYTES_POST = []byte("(\xacY\xc1\xf4\xa7\xce\xac\xdd\xaaO)\n")

	CROSS_POST_PREFIX = []byte("[\xc2\xe0\xbf\xfd]")
)
