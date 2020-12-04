package ptttype

type SortByCategory int

const (
	SORT_BY_ID SortByCategory = iota
	SORT_BY_CLASS
	SORT_BY_STAT
	SORT_BY_IDLE
	SORT_BY_FROM
	SORT_BY_FIVE
	SORT_BY_SEX

	SORT_BY_MAX SortByCategory = 9
)

func (d SortByCategory) String() string {
	switch d {
	case SORT_BY_ID:
		return "id"
	case SORT_BY_CLASS:
		return "class"
	case SORT_BY_STAT:
		return "stat"
	case SORT_BY_IDLE:
		return "idle"
	case SORT_BY_FROM:
		return "from"
	case SORT_BY_FIVE:
		return "five"
	case SORT_BY_SEX:
		return "gender"
	default:
		return "unknown"
	}
}
