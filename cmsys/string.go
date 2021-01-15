package cmsys

import (
	"bytes"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
)

func StringHashWithHashBits(theBytes []byte) uint32 {
	return StringHash(theBytes) % (1 << ptttype.HASH_BITS)
}

func StringHash(theBytes []byte) uint32 {
	return fnv1a32StrCase(theBytes, FNV1_32_INIT)
}

func StripBlank(theBytes []byte) []byte {
	theIdx := bytes.Index(theBytes, []byte{' '})
	if theIdx == -1 {
		return theBytes
	}

	return theBytes[:theIdx]
}

//StringNoneBig5
//
//https://github.com/ptt/pttbbs/issues/94
//str is modified.
func StripNoneBig5(str_out []byte) (sanitizedStr []byte) {
	theLen := 0
	for idx := 0; idx < len(str_out) && str_out[idx] != 0; idx++ {
		if 32 <= str_out[idx] && str_out[idx] < 128 {
			str_out[theLen] = str_out[idx]
			theLen++
		} else if str_out[idx] == 255 {
			if idx+1 < len(str_out) {
				if 251 <= str_out[idx+1] && str_out[idx+1] <= 254 {
					idx++
					//XXX this looks strange,
					//because we've already got idx++ in the end of for-loop,
					//but this setup is to be compatible with c-pttbbs.
					if idx+1 < len(str_out) && str_out[idx+1] != 0 {
						idx++
					}
				}
			}
			continue
		} else if (str_out[idx] & 0x80) != 0 {
			if idx+1 < len(str_out) {
				if 0x40 <= str_out[idx+1] && str_out[idx+1] <= 0x7e ||
					0xa1 <= str_out[idx+1] && str_out[idx+1] <= 0xfe {
					str_out[theLen] = str_out[idx]
					theLen++
					str_out[theLen] = str_out[idx+1]
					theLen++
					idx++
				}
			}
		}
	}

	if theLen < len(str_out) {
		str_out[theLen] = 0
	}

	return str_out[:theLen]
}

func StripAnsi(src []byte, flag StripAnsiFlag) (dst []byte) {
	dst = make([]byte, len(src))

	idxDst := 0
	for idxSrc := 0; idxSrc < len(src); idxSrc++ {
		each := src[idxSrc]
		if each == 0 {
			break
		}

		if each != ansi.ESC_CHR {
			if idxDst < len(dst) {
				dst[idxDst] = each
				idxDst++
			}
			continue
		}

		if idxSrc == len(src)-1 { //the last char
			break
		}

		idxP := idxSrc + 1
		p := src[idxP]
		if p != '[' { //exception
			//XXX we would like to skip the char following \x1e
			idxSrc += 1
			if src[idxSrc] == 0 {
				break
			}
			continue
		}

		for idxP = idxP + 1; idxP < len(src) && isEscapeParam(src[idxP]); idxP = idxP + 1 {
		}

		p = src[idxP]
		if (flag == STRIP_ANSI_NO_RELOAD && isEscapeCommand(src[idxP])) || (flag == STRIP_ANSI_ONLY_COLOR && src[idxP] == 'm') {
			theLen := idxP - idxSrc + 1 //len dst is same as src and idxDst is < idxSrc. we don't need to worry buffer overflow in dst for now.
			copy(dst[idxDst:(idxDst+theLen)], src[idxSrc:(idxSrc+theLen)])
			idxDst += theLen
		}

		idxSrc = idxP
		if idxSrc == len(src) || src[idxSrc] == 0 {
			break
		}
	}

	return dst[:idxDst]
}

func isEscapeParam(x byte) bool {
	return ESCAPE_FLAG[x]&1 != 0
}

func isEscapeCommand(x byte) bool {
	return ESCAPE_FLAG[x]&2 != 0
}
