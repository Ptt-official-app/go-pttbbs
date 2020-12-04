package ptttype

func ValidUSHMEntry(x int) bool {
	return x >= 0 && x < USHM_SIZE
}
