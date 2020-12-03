package types

func Isalpha(c byte) bool {
	if c >= 'A' && c <= 'Z' {
		return true
	}

	if c >= 'a' && c <= 'z' {
		return true
	}

	return false
}

func Isnumber(c byte) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func Isalnum(c byte) bool {
	return Isalpha(c) || Isnumber(c)
}
