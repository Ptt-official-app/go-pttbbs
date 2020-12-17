package ptttype

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

//We have 3 different ids for user:
//	UserID_t: (username)
//	Uid: (int32) (uid starting from 1)
//  UidInStore: (int32) (Uid - 1)
type UserID_t [IDLEN + 1]byte
type Uid int32
type UidInStore int32
type RealName_t [REALNAMESZ]byte
type Nickname_t [NICKNAMESZ]byte
type Passwd_t [PASSLEN]byte
type IPv4_t [IPV4LEN + 1]byte
type Email_t [EMAILSZ]byte
type Address_t [ADDRESSSZ]byte
type Reg_t [REGLEN + 1]byte
type Career_t [CAREERSZ]byte
type Phone_t [PHONESZ]byte

type ChatID_t [11]byte

type From_t [27]byte

type Date_t [6]byte

//We have 3 different ids for board:
//  BoardID_t: (brdname)
//  Bid: (int32) (bid starting from 1)
//  BidInStore (int32) (Bid - 1)
type BoardID_t [IDLEN + 1]byte
type Bid int32
type BidInStore int32
type BoardTitle_t [BTLEN + 1]byte
type BM_t [IDLEN*3 + 3]byte /* BMs' userid, token '/' */

type Aid int32
type AidInStore int32
type Aidu uint64 /* ptt-aidu */
type Aidc [8]byte

type Filename_t [FNLEN]byte
type Subject_t [STRLEN]byte
type RCPT_t [RCPTSZ]byte

type Owner_t [IDLEN + 2]byte //user-id[.]

type Title_t [TTLEN + 1]byte

type Gid int32

type SortIdx int

var (
	EMPTY_USER_ID     = UserID_t{}
	EMPTY_BOARD_ID    = BoardID_t{}
	EMPTY_BOARD_TITLE = BoardTitle_t{}
)

const USER_ID_SZ = unsafe.Sizeof(EMPTY_USER_ID)
const BOARD_ID_SZ = unsafe.Sizeof(EMPTY_BOARD_ID)
const BOARD_TITLE_SZ = unsafe.Sizeof(EMPTY_BOARD_TITLE)
const UID_IN_STORE_SZ = unsafe.Sizeof(UidInStore(0))
const UID_SZ = unsafe.Sizeof(Uid(0))
const BID_IN_STORE_SZ = unsafe.Sizeof(BidInStore(0))
const BID_SZ = unsafe.Sizeof(Bid(0))

func (u UidInStore) ToUid() Uid {
	return Uid(u + 1)
}

func (u Uid) ToUidInStore() UidInStore {
	return UidInStore(u - 1)
}

func (u Uid) String() string {
	return strconv.Itoa(int(u))
}

func (u Uid) IsValid() bool {
	return u >= 1 && u <= MAX_USERS
}

func (b Bid) IsValid() bool {
	return b >= 1 && b <= MAX_BOARD
}

func (b BidInStore) ToBid() Bid {
	return Bid(b + 1)
}

func (b Bid) ToBidInStore() BidInStore {
	return BidInStore(b - 1)
}

func (b Bid) String() string {
	return strconv.Itoa(int(b))
}

func (a Aid) ToAidInStore() AidInStore {
	return AidInStore(a - 1)
}

func (a Aid) IsValid() bool {
	return a >= 1
}

func (a AidInStore) ToAid() Aid {
	return Aid(a + 1)
}

func (s SortIdx) String() string {
	if s < 0 {
		return ""
	}
	return strconv.Itoa(int(s))
}

func ToSortIdx(str string) (SortIdx, error) {
	if str == "" {
		return 0, nil
	}

	idx, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}

	return SortIdx(idx), nil
}

//valid UserID
//
//https://github.com/ptt/pttbbs/blob/master/common/bbs/names.c
func (u *UserID_t) IsValid() bool {
	if u == nil {
		return false
	}

	theLen := types.Cstrlen(u[:])
	if theLen < 2 || theLen > IDLEN {
		return false
	}

	if !types.Isalpha(u[0]) {
		return false
	}

	for idx, c := range u {
		if idx == theLen {
			break
		}

		if !types.Isalnum(c) {
			return false
		}
	}
	return true
}

//Valid BoardID
//
//https://github.com/ptt/pttbbs/blob/master/common/bbs/string.c#L21
func (b *BoardID_t) IsValid() bool {
	lenB := types.Cstrlen(b[:])
	if lenB < 2 || lenB > IDLEN {
		return false
	}

	ch := b[0]
	if !types.Isalpha(ch) {
		return false
	}
	for idx := 1; idx < lenB; idx++ {
		if !types.Isalnum(ch) && ch != '_' && ch != '-' && ch != '.' {
			return false
		}
	}

	return true
}

//RealTitle
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) RealTitle() []byte {
	return types.CstrToBytes(t[7:41])
}

//BoardClass
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) BoardClass() []byte {
	result := t[:5]
	if result[4] == ' ' {
		result = result[:4]
	}
	return result
}

//BoardType
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) BoardType() []byte {
	return t[5:7]
}

//ToBMs
//
//We would like to have a better method
//(We don't need to worry about this once we move everything to the db.)
func (bm *BM_t) ToBMs() []*UserID_t {
	bmBytes := types.CstrToBytes(bm[:])
	theList := bytes.Split(bmBytes, []byte{'/'})
	bms := make([]*UserID_t, 0, len(theList))
	for _, each := range theList {
		if len(each) == 0 || each[0] == 0 {
			continue
		}
		eachUser := &UserID_t{}
		copy(eachUser[:], each[:])
		bms = append(bms, eachUser)
	}

	return bms
}

//////////
//Filename
//////////

func (f *Filename_t) Type() RecordType {
	switch f[0] {
	case 'M':
		return RECORD_TYPE_M
	default:
		return RECORD_TYPE_G
	}
}
func (f *Filename_t) CreateTime() (types.Time4, error) {
	createTime_i, err := strconv.Atoi(string(f[2:12]))
	if err != nil {
		return 0, err
	}
	return types.Time4(createTime_i), nil
}

const FILENAME_OFFSET_POSTFIX = 15

func (f *Filename_t) Postfix() []byte {
	return f[FILENAME_OFFSET_POSTFIX:(FILENAME_OFFSET_POSTFIX + 3)]
}

func (f *Filename_t) ToAidu() Aidu {
	if f[1] != '.' {
		return 0
	}
	if f[12] != '.' {
		return 0
	}
	if f[14] != '.' {
		return 0
	}
	theType := f.Type()

	createTime, _ := f.CreateTime()

	postfixStr := string(f.Postfix())
	postfix_u64, err := strconv.ParseUint(postfixStr, 16, 12)
	if err != nil {
		return 0
	}

	aidu_u64 := (uint64(theType&0xf) << 44) + uint64(createTime)<<12 + postfix_u64
	return Aidu(aidu_u64)
}

//Type
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L82
func (a Aidu) Type() RecordType {
	theType := (a >> 44) & 0xf
	switch theType {
	case 0:
		return RECORD_TYPE_M
	default:
		return RECORD_TYPE_G
	}
}

//Time
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L83
func (a Aidu) Time() types.Time4 {
	return types.Time4((a >> 12) & 0xffffffff)
}

func (a Aidu) Postfix() uint16 {
	return uint16(a) & 0xfff
}

//FN
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L80
func (a Aidu) ToFN() *Filename_t {
	theType := a.Type()
	theTime := a.Time()
	postfix := a.Postfix()

	fnStr := fmt.Sprintf("%v.%d.A.%03x", theType, theTime, postfix)
	fn := &Filename_t{}
	copy(fn[:], fnStr)
	return fn
}

const encodeAidc = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
const lenEncodeAidc = uint64(len(encodeAidc))

func (a Aidu) ToAidc() *Aidc {
	aidc := &Aidc{}

	aidu := uint64(a)
	for idx := len(aidc) - 1; idx >= 0; idx-- {
		v := aidu % lenEncodeAidc
		aidu = aidu / lenEncodeAidc
		aidc[idx] = encodeAidc[v]
	}
	return aidc
}

var decodeAidcTable = [128]Aidu{
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // , !, ", #, $, %, &, '
	0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x00, 0x00, //(, ), *, +, ,, -, ., /
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, //0, 1, 2, 3, 4, 5, 6, 7
	0x08, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //8, 9, :, ;, <, =, >, ?
	0x00, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, //@, A, B, C, D, E, F, G
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, //H, I, J, K, L, M, N, O
	0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, //P, Q, R, S, T, U, V, W
	0x21, 0x22, 0x23, 0x00, 0x00, 0x00, 0x00, 0x3f, //X, Y, Z, [, \, ], ^, _
	0x00, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, //`, a, b, c, d, e, f, g
	0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, //h, i, j, k, l, m, n, o
	0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, //p, q, r, s, t, u, v, w
	0x3b, 0x3c, 0x3d, 0x00, 0x00, 0x00, 0x00, 0x00, //x, y, z, {, |, }, ~, [del]
}

func (a *Aidc) ToAidu() (aidu Aidu) {
	aidu = 0
	for _, each := range a {
		if each == 0 {
			break
		}
		if each == '@' {
			break
		}

		v := decodeAidcTable[each]
		aidu <<= 6
		aidu |= v
	}
	return aidu
}

//////////
//Owner
//////////
func (o *Owner_t) ToUserID() *UserID_t {
	userID := &UserID_t{}
	oBytes := types.CstrToBytes(o[:])
	if !types.Isalnum(oBytes[len(oBytes)-1]) {
		oBytes = oBytes[:len(oBytes)-1]
	}
	copy(userID[:], oBytes[:])
	return userID
}
