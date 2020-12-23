package ptttype

type RestrictReason uint8

const (
	_ RestrictReason = iota
	RESTRICT_REASON_FORBIDDEN
	RESTRICT_REASON_HIDDEN
)

func (r RestrictReason) String() string {
	switch r {
	case RESTRICT_REASON_FORBIDDEN:
		return "forbidden"
	case RESTRICT_REASON_HIDDEN:
		return "hidden"
	default:
		return ""
	}
}
