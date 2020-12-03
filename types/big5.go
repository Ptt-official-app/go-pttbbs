package types

import (
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//Utf8ToBig5
//
//We have several fields with pure big5 Cstr.
//use DBCUToDBCS if it's with color-decorators.
func Utf8ToBig5(input string) Cstr {
	utf8ToBig5 := traditionalchinese.Big5.NewEncoder()
	big5, _, _ := transform.Bytes(utf8ToBig5, []byte(input))
	return big5
}

//Big5ToUtf8
//
//We have several fields with pure big5 Cstr.
//use DBCSToDBCU if it's with color-decorators.
func Big5ToUtf8(input Cstr) string {
	big5ToUTF8 := traditionalchinese.Big5.NewDecoder()
	inputBytes := CstrToBytes(input)
	utf8, _, _ := transform.String(big5ToUTF8, string(inputBytes))
	return utf8
}
