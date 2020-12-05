package ptttype

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/types"
)

type FileHeaderRaw struct {
	Filename  Filename_t  /* M.1120582370.A.1EA [19+1], create time */
	Modified  types.Time4 /* last modified time */
	Pad       byte        /* padding, not used */
	Recommend int8        /* important level */
	Owner     Owner_t     /* uid[.] */
	Date      Date_t      /* [02/02] or space(5) */
	Title     Title_t
	/* TODO this multi is a mess now. */
	Pad2  byte
	Multi [4]byte //union as either money (int) or anon_uid (int) or vote_limits (4 unsigned char) or refer (2 unsigned int)

	//union {
	///* TODO: MOVE money to outside multi!!!!!! */
	//int money;
	//int anon_uid;
	///* different order to match alignment */
	//struct {
	//    unsigned char posts;
	//    unsigned char logins;
	//    unsigned char regtime;
	//    unsigned char badpost;
	//} vote_limits;
	//struct {
	//    /* is this ordering correct? */
	//    unsigned int ref:31;
	//    unsigned int flag:1;
	//} refer;
	//}	    multi;		    /* rocker: if bit32 on ==> reference */
	/* XXX dirty, split into flag and money if money of each file is less than 16bit? */
	Filemode uint8 /* must be last field @ boards.c */
	Pad3     [3]byte
}

var emptyFileHeaderRaw = FileHeaderRaw{}

const FILE_HEADER_RAW_SZ = unsafe.Sizeof(emptyFileHeaderRaw)

//XXX need to ensure Multi.
type VoteLimits struct {
	Post    uint8
	Logins  uint8
	RegTime uint8
	Badpost uint8
}

//XXX need to ensure FileRefer.
type FileRefer uint32

func (f FileRefer) Ref() uint32 {
	return uint32(f & 0x7fffffff)
}

func (f FileRefer) Flag() uint8 {
	return uint8(f >> 31)
}

func (f *FileHeaderRaw) Money() (money int32) {
	buf := bytes.NewBuffer(f.Multi[:4])
	_ = binary.Read(buf, binary.LittleEndian, &money)
	return money
}

func (f *FileHeaderRaw) AnonUID() (anonUID int32) {
	buf := bytes.NewBuffer(f.Multi[:4])
	_ = binary.Read(buf, binary.LittleEndian, &anonUID)
	return anonUID
}

func (f *FileHeaderRaw) VoteLimits() *VoteLimits {
	return &VoteLimits{f.Multi[0], f.Multi[1], f.Multi[2], f.Multi[3]}
}

func (f *FileHeaderRaw) VoteLimitPosts() uint8 {
	return f.Multi[0]
}

func (f *FileHeaderRaw) VoteLimitLogins() uint8 {
	return f.Multi[1]
}

func (f *FileHeaderRaw) VoteLimitRegTime() uint8 {
	return f.Multi[2]
}

func (f *FileHeaderRaw) VoteLimitBadpost() uint8 {
	return f.Multi[3]
}
