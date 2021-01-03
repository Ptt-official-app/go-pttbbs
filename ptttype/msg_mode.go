package ptttype

type MsgMode int32

const (
	MSGMODE_TALK      MsgMode = 0
	MSGMODE_WRITE     MsgMode = 1
	MSGMODE_FROMANGEL MsgMode = 2
	MSGMODE_TOANGEL   MsgMode = 3
	MSGMODE_ALOHA     MsgMode = 4
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
