package ptttype

type RecordType uint8

const (
	RECORD_TYPE_M RecordType = iota
	RECORD_TYPE_G            // not RECORD_TYPE_M will be TYPE_G
)

func (r RecordType) String() string {
	switch r {
	case RECORD_TYPE_M:
		return "M"
	default:
		return "G"
	}
}
