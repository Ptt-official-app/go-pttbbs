package ptttype

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

// We have 3 different ids for user:
//   - UserID_t: (username)
//   - UID: (int32) (uid starting from 1)
//   - UIDInStore: (int32) (UID - 1)
type (
	UserID_t   [IDLEN + 1]byte
	UID        int32
	UIDInStore int32
	RealName_t [REALNAMESZ]byte
	Nickname_t [NICKNAMESZ]byte
	Passwd_t   [PASSLEN]byte
	IPv4_t     [IPV4LEN + 1]byte
	Email_t    [EMAILSZ]byte
	Address_t  [ADDRESSSZ]byte
	Reg_t      [REGLEN + 1]byte
	Career_t   [CAREERSZ]byte
	Phone_t    [PHONESZ]byte
)

type UtmpID int32 // starting from 0, idx in Shm.Raw.UInfo

type ChatID_t [11]byte

type From_t [27]byte

type Date_t [6]byte

// We have 3 different ids for board:
//   - BoardID_t: (brdname)
//   - Bid: (int32) (bid starting from 1)
//   - BidInStore (int32) (Bid - 1)
type (
	BoardID_t    [IDLEN + 1]byte
	Bid          int32
	BidInStore   int32
	BoardTitle_t [BTLEN + 1]byte
	BM_t         [IDLEN*3 + 3]byte /* BMs' userid, token '/' */)

type (
	Aid        int32
	AidInStore int32
	Aidu       uint64 /* ptt-aidu */
	Aidc       [8]byte
)

type (
	Filename_t [FNLEN]byte
	Subject_t  [STRLEN]byte
	RCPT_t     [RCPTSZ]byte
)

type Owner_t [IDLEN + 2]byte // user-id[.]

type Title_t [TTLEN + 1]byte

type (
	SortIdx        int
	SortIdxInStore int
)

var (
	EMPTY_USER_ID     = UserID_t{}
	EMPTY_BOARD_ID    = BoardID_t{}
	EMPTY_BOARD_TITLE = BoardTitle_t{}
	EMPTY_EMAIL       = Email_t{}
	EMPTY_BM          = BM_t{}
	EMPTY_AIDC        = Aidc{}
)

const (
	USER_ID_SZ      = unsafe.Sizeof(EMPTY_USER_ID)
	BOARD_ID_SZ     = unsafe.Sizeof(EMPTY_BOARD_ID)
	BOARD_TITLE_SZ  = unsafe.Sizeof(EMPTY_BOARD_TITLE)
	UID_IN_STORE_SZ = unsafe.Sizeof(UIDInStore(0))
	UTMP_ID_SZ      = unsafe.Sizeof(UtmpID(0))
	UID_SZ          = unsafe.Sizeof(UID(0))
	BID_IN_STORE_SZ = unsafe.Sizeof(BidInStore(0))
	BID_SZ          = unsafe.Sizeof(Bid(0))
	EMAIL_SZ        = unsafe.Sizeof(EMPTY_EMAIL)
	BM_SZ           = unsafe.Sizeof(EMPTY_BM)
	AIDC_SZ         = unsafe.Sizeof(EMPTY_AIDC)
)

func (u UIDInStore) ToUID() UID {
	return UID(u + 1)
}

func (u UID) ToUIDInStore() UIDInStore {
	return UIDInStore(u - 1)
}

func (u UID) String() string {
	return strconv.Itoa(int(u))
}

func (u UID) IsValid() bool {
	return u >= 1 && u <= MAX_USERS
}

func (u UID) ToPid() (p types.Pid_t) {
	return types.Pid_t(u) + types.DEFAULT_PID_MAX
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

func (s SortIdx) ToSortIdxInStore() SortIdxInStore {
	return SortIdxInStore(s - 1)
}

func (s SortIdx) IsValid() bool {
	return s >= 1
}

func (s SortIdxInStore) ToSortIdx() SortIdx {
	return SortIdx(s + 1)
}

// valid UserID
//
// https://github.com/ptt/pttbbs/blob/master/common/bbs/names.c
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

func (u *UserID_t) CopyFrom(uBytes []byte) {
	copy(u[:], uBytes)
	if len(uBytes) < int(USER_ID_SZ) {
		u[len(uBytes)] = 0
	}
}

// IsGuest
// guest as reserved user-account.
func (u *UserID_t) IsGuest() bool {
	return u[0] == 'g' && u[1] == 'u' && u[2] == 'e' && u[3] == 's' && u[4] == 't' && u[5] == 0
}

func ToBoardID(boardIDBytes []byte) (boardID *BoardID_t) {
	boardID = &BoardID_t{}
	copy(boardID[:], boardIDBytes)

	return boardID
}

// Valid BoardID
//
// https://github.com/ptt/pttbbs/blob/master/common/bbs/string.c#L21
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

// RealTitle
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) RealTitle() []byte {
	return types.CstrToBytes(t[7:41])
}

// BoardClass
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) BoardClass() []byte {
	result := t[:5]
	if result[4] == ' ' { // no ' ' in big5-encoding, we can safely remove ' '.
		result = result[:4]
	}
	return result
}

// BoardType
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) BoardType() []byte {
	return t[5:7]
}

// ToBMs
//
// We would like to have a better method
// (We don't need to worry about this once we move everything to the db.)
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

func (f *Filename_t) String() string {
	return types.CstrToString(f[:])
}

// Eq
//
// It's possible that the timestamp
// Compare only with the timestamp and the rnd.
func (f *Filename_t) Eq(f2 *Filename_t) bool {
	return types.Cstrcmp(f[2:], f2[2:]) == 0
}

func (f *Filename_t) IsDeleted() bool {
	return bytes.Equal(f[:FN_SAFEDEL_PREFIX_LEN], FN_SAFEDEL_b)
}

func (f *Filename_t) Basename() string {
	if f.IsDeleted() {
		return "M." + string(f[FN_SAFEDEL_PREFIX_LEN:])
	}

	return types.CstrToString(f[:])
}

func (f *Filename_t) DeletedName() string {
	return FN_SAFEDEL + types.CstrToString(f[FN_SAFEDEL_PREFIX_LEN:])
}

// Type
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L82
func (a Aidu) Type() RecordType {
	theType := (a >> 44) & 0xf
	switch theType {
	case 0:
		return RECORD_TYPE_M
	default:
		return RECORD_TYPE_G
	}
}

// Time
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L83
func (a Aidu) Time() types.Time4 {
	return types.Time4((a >> 12) & 0xffffffff)
}

func (a Aidu) Postfix() uint16 {
	return uint16(a) & 0xfff
}

// FN
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/aids.c#L80
func (a Aidu) ToFN() *Filename_t {
	theType := a.Type()
	theTime := a.Time()
	postfix := a.Postfix()

	fnStr := fmt.Sprintf("%v.%d.A.%03X", theType, theTime, postfix)
	fn := &Filename_t{}
	copy(fn[:], fnStr)
	return fn
}

const (
	encodeAidc    = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	lenEncodeAidc = uint64(len(encodeAidc))
)

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
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, // 0, 1, 2, 3, 4, 5, 6, 7
	0x08, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 8, 9, :, ;, <, =, >, ?
	0x00, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, //@, A, B, C, D, E, F, G
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, // H, I, J, K, L, M, N, O
	0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // P, Q, R, S, T, U, V, W
	0x21, 0x22, 0x23, 0x00, 0x00, 0x00, 0x00, 0x3f, // X, Y, Z, [, \, ], ^, _
	0x00, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, //`, a, b, c, d, e, f, g
	0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, // h, i, j, k, l, m, n, o
	0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, // p, q, r, s, t, u, v, w
	0x3b, 0x3c, 0x3d, 0x00, 0x00, 0x00, 0x00, 0x00, // x, y, z, {, |, }, ~, [del]
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

// ToUserID
func (o *Owner_t) ToUserID() *UserID_t {
	userID := &UserID_t{}
	oBytes := types.CstrToBytes(o[:])
	if !types.Isalnum(oBytes[len(oBytes)-1]) {
		oBytes = oBytes[:len(oBytes)-1]
	}
	copy(userID[:], oBytes[:])
	return userID
}

// https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L653
func (o *Owner_t) IsCorpse() bool {
	return o[0] == '-' && o[1] == 0
}

// NewBM
//
// called only in cache.ParseBMList.
// Already verified in cache.ParseBMList.
// no need to worry that userIDs exceeds BM_t
func NewBM(userIDs []*UserID_t) (bms *BM_t) {
	bms = &BM_t{}
	bmsBytes := bms[:]
	p_bmsBytes := bmsBytes
	for idx, each := range userIDs {
		if idx > 0 {
			p_bmsBytes[0] = '/'
			p_bmsBytes = p_bmsBytes[1:]
		}
		userIDBytes := types.CstrToBytes(each[:])
		copy(p_bmsBytes[:], userIDBytes)
		p_bmsBytes = p_bmsBytes[len(userIDBytes):]
	}

	return bms
}
