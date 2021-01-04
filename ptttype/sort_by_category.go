package ptttype

type SortByCategory int

const (
	SORT_BY_ID    SortByCategory = 0
	SORT_BY_CLASS SortByCategory = 1
	SORT_BY_STAT  SortByCategory = 2
	SORT_BY_IDLE  SortByCategory = 3
	SORT_BY_FROM  SortByCategory = 4
	SORT_BY_FIVE  SortByCategory = 5
	SORT_BY_SEX   SortByCategory = 6

	SORT_BY_UID SortByCategory = 7

	SORT_BY_PID SortByCategory = 8

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
