package types

import (
	"bytes"
)

//Cstr
//
//[]byte with C String property in that \0 is considered as the end of the bytes/string.
//It is used to convert from fixed-length bytes to string or []byte with no \0.
//
//Naming Cstr instead of CString is to avoid confusion with C.CString
//(C.CString is from string, and should be compatible with string, not with []byte)
//(We also have str(len/cpy/cmp) functions in C)
//
//See tests for more examples of how to use fixed-bytes with Cstr to get no-\0 string / []byte
type Cstr []byte

func Cstrlen(cstr Cstr) int {
	theLen := bytes.IndexByte(cstr, 0x00)
	if theLen == -1 {
		return len(cstr)
	}

	return theLen
}

//CstrToString
//
//Only the bytes until \0 when converting to string.
//See tests for more examples.
//
//Params
//	cstr
//
//Return
//	string: string
func CstrToString(cstr Cstr) string {
	theBytes := CstrToBytes(cstr)
	return string(theBytes)
}

//CstrToBytes
//
//Only the bytes until \0.
//See tests for more examples.
//
//Params
//	cstr
//
//Return
//	[]byte: bytes
func CstrToBytes(cstr Cstr) []byte {
	theLen := Cstrlen(cstr)
	if theLen == 0 {
		return nil
	}
	return cstr[:theLen]
}

//Cstrcmp
//The reason that we don't directly do C.strcmp is because cstr1 or cstr2 may not have \x00 in the end.
func Cstrcmp(cstr1 Cstr, cstr2 Cstr) int {
	// iterate through cstr1
	len1 := len(cstr1)
	len2 := len(cstr2)
	var each2 byte
	for idx, each := range cstr1 {
		if each == 0 { // reach the end of cstr1
			len1 = idx
			if idx < len2 && cstr2[idx] == 0 { // possibly len2 reaches the end as well.
				len2 = idx
			}
			break
		}

		// if current-idx > len2 (abcde vs. abc)
		if idx >= len2 {
			return int(each)
		}
		each2 = cstr2[idx]
		if each != each2 { //it's ok if each2 == 0, because because it means that each - each2 > 0
			return int(each) - int(each2)
		}
	}

	// until now, we know that cstr1 and cstr2 are the same with len1.
	if len1 < len2 {
		return -int(cstr2[len1])
	}

	return 0
}

func Cstrcasecmp(cstr1 Cstr, cstr2 Cstr) int {
	cstr1Lower := CstrTolower(cstr1)
	cstr2Lower := CstrTolower(cstr2)

	return Cstrcmp(cstr1Lower, cstr2Lower)
}

func Cstrstr(cstr Cstr, substr Cstr) int {
	theIdx := bytes.Index(cstr, substr)
	theLen := Cstrlen(cstr)
	if theIdx < 0 || theIdx >= theLen {
		return -1
	}

	return theIdx
}

func Cstrcasestr(cstr Cstr, substr Cstr) int {
	cstrLower := CstrTolower(cstr)
	substrLower := CstrTolower(substr)

	return Cstrstr(cstrLower, substrLower)
}

func CstrTolower(cstr Cstr) Cstr {
	cstrLower := make(Cstr, len(cstr))
	for idx, each := range cstr {
		cstrLower[idx] = CcharTolower(each)
	}

	return cstrLower
}

func CcharTolower(ch byte) byte {
	if ch >= 'A' && ch <= 'Z' {
		return ch + 'a' - 'A'
	}

	return ch
}

func CstrToupper(cstr Cstr) Cstr {
	cstrUpper := make(Cstr, len(cstr))
	for idx, each := range cstr {
		cstrUpper[idx] = CcharToupper(each)
	}

	return cstrUpper
}

func CcharToupper(ch byte) byte {
	if ch >= 'a' && ch <= 'z' {
		return ch - 32
	}
	return ch
}
