package ptttype

import "unsafe"

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

var (
	EMPTY_USER_ID     = UserID_t{}
	EMPTY_BOARD_ID    = BoardID_t{}
	EMPTY_BOARD_TITLE = BoardTitle_t{}
)

const USER_ID_SZ = unsafe.Sizeof(EMPTY_USER_ID)
const BOARD_ID_SZ = unsafe.Sizeof(EMPTY_BOARD_ID)
const BOARD_TITLE_SZ = unsafe.Sizeof(EMPTY_BOARD_TITLE)
const UID_IN_STORE_SZ = unsafe.Sizeof(UidInStore(0))

func (u UidInStore) ToUid() Uid {
	return Uid(u + 1)
}

func (u Uid) ToUidInStore() UidInStore {
	return UidInStore(u - 1)
}

func (b BidInStore) ToBid() Bid {
	return Bid(b + 1)
}

func (b Bid) ToBidInStore() BidInStore {
	return BidInStore(b - 1)
}
