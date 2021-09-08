package cmsys

import "github.com/Ptt-official-app/go-pttbbs/types"

// https://github.com/ptt/pttbbs/blob/master/include/fnv_hash.h
// commit: 6bdd36898bde207683a441cdffe2981e95de5b20

func fnv32Bytes(theBytes []byte, hval Fnv32_t) Fnv32_t {
	for _, each := range theBytes {
		hval *= FNV_32_PRIME
		hval ^= Fnv32_t(each)
	}

	return hval
}

func fnv1a32Bytes(theBytes []byte, hval Fnv32_t) Fnv32_t {
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		hval ^= Fnv32_t(each)
		hval *= FNV_32_PRIME
	}

	return hval
}

func fnv1a32StrCase(theBytes []byte, hval Fnv32_t) Fnv32_t {
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		hval ^= Fnv32_t(types.CcharToupper(each))
		hval *= FNV_32_PRIME
	}

	return hval
}

func fnv1a32DBCSCase(theBytes []byte, hval Fnv32_t) Fnv32_t {
	isDBCS := false
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		if isDBCS {
			// 2nd DBCS
			isDBCS = false
		} else {
			if each < 0x80 {
				each = types.CcharToupper(each)
			} else {
				isDBCS = true
			}
		}
		hval ^= Fnv32_t(each)
		hval *= FNV_32_PRIME
	}

	return hval
}

//////////
//64bits
//////////

func Fnv64Buf(buf []byte, theLen int, hval Fnv64_t) (newHVal Fnv64_t) {
	for _, each := range buf {
		hval *= FNV_64_PRIME
		hval ^= Fnv64_t(each)

		theLen--
		if theLen == 0 {
			break
		}
	}
	return hval
}

func fnv64Bytes(theBytes []byte, hval Fnv64_t) Fnv64_t {
	for _, each := range theBytes {
		hval *= FNV_64_PRIME
		hval ^= Fnv64_t(each)
	}

	return hval
}

func fnv1a64Bytes(theBytes []byte, hval Fnv64_t) Fnv64_t {
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		hval ^= Fnv64_t(each)
		hval *= FNV_64_PRIME
	}

	return hval
}

func fnv1a64StrCase(theBytes []byte, hval Fnv64_t) Fnv64_t {
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		hval ^= Fnv64_t(types.CcharToupper(each))
		hval *= FNV_64_PRIME
	}

	return hval
}

func fnv1a64DBCSCase(theBytes []byte, hval Fnv64_t) Fnv64_t {
	isDBCS := false
	for _, each := range theBytes {
		if each == 0 {
			break
		}
		if isDBCS {
			// 2nd DBCS
			isDBCS = false
		} else {
			if each < 0x80 {
				each = types.CcharToupper(each)
			} else {
				isDBCS = true
			}
		}
		hval ^= Fnv64_t(each)
		hval *= FNV_64_PRIME
	}

	return hval
}

func fnv1aByte(theByte byte, hval Fnv32_t) Fnv32_t {
	hval ^= Fnv32_t(theByte)
	hval *= FNV_32_PRIME
	return hval
}
