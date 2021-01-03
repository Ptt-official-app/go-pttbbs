package ptttype

func ValidUSHMEntry(x UtmpID) bool {
	return x >= 0 && x < USHM_SIZE
}
