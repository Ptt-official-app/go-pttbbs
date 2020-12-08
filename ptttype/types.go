package ptttype

import (
	"bytes"
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

func (b BidInStore) ToBid() Bid {
	return Bid(b + 1)
}

func (b Bid) ToBidInStore() BidInStore {
	return BidInStore(b - 1)
}

func (b Bid) String() string {
	return strconv.Itoa(int(b))
}

func (s SortIdx) String() string {
	if s <= 0 {
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

//RealTitle
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) RealTitle() []byte {
	return t[7:41]
}

//BoardClass
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/board.c#L1517
func (t *BoardTitle_t) BoardClass() []byte {
	return t[:4]
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
