package ptttype

type MsgMode uint8

const (
	MSGMODE_TALK MsgMode = iota
	MSGMODE_WRITE
	MSGMODE_FROMANGEL
	MSGMODE_TOANGEL
	MSGMODE_ALOHA
)

func (m MsgMode) String() string {
	switch m {
	case MSGMODE_TALK:
		return "talk"
	case MSGMODE_WRITE:
		return "write"
	case MSGMODE_FROMANGEL:
		return "from-angel"
	case MSGMODE_TOANGEL:
		return "to-angel"
	default:
		return "[unknown]"
	}
}
