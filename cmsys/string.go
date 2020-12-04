package cmsys

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

func StringHashWithHashBits(theBytes []byte) uint32 {
	return StringHash(theBytes) % (1 << ptttype.HASH_BITS)
}

func StringHash(theBytes []byte) uint32 {
	return fnv1a32StrCase(theBytes, FNV1_32_INIT)
}
