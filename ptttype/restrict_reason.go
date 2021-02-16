package ptttype

type RestrictReason uint8

const (
	RESTRICT_REASON_NONE          RestrictReason = 0
	RESTRICT_REASON_FORBIDDEN     RestrictReason = 1
	RESTRICT_REASON_HIDDEN        RestrictReason = 2
	RESTRICT_REASON_NUMLOGIN_DAYS RestrictReason = 3
	RESTRICT_REASON_BADPOST       RestrictReason = 4
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
