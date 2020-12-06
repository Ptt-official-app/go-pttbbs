package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

// https://github.com/ptt/pttbbs/blob/master/include/pttstruct.h
type Userec struct {
	Version  uint32
	Userid   string
	Realname string
	Nickname string
	Passwd   string
	Pad1     uint8

	Uflag        uint32
	_unused1     uint32
	Userlevel    uint32
	Numlogindays uint32
	Numposts     uint32
	Firstlogin   uint32
	Lastlogin    uint32
	Lasthost     string
	// TODO
}

func NewUserecFromRaw(userecRaw *ptttype.UserecRaw) *Userec {
	user := &Userec{}
	user.Version = userecRaw.Version
	user.Userid = types.CstrToString(userecRaw.UserID[:])
	user.Realname = types.Big5ToUtf8(types.CstrToBytes(userecRaw.RealName[:]))
	user.Nickname = types.Big5ToUtf8(types.CstrToBytes(userecRaw.Nickname[:]))
	user.Passwd = types.CstrToString(userecRaw.PasswdHash[:])
	user.Pad1 = userecRaw.Pad1

	user.Uflag = uint32(userecRaw.UFlag)
	user._unused1 = userecRaw.Unused1
	user.Userlevel = uint32(userecRaw.UserLevel)
	user.Numlogindays = userecRaw.NumLoginDays
	user.Numposts = userecRaw.NumPosts
	user.Firstlogin = uint32(userecRaw.FirstLogin)
	user.Lastlogin = uint32(userecRaw.LastLogin)
	user.Lasthost = types.CstrToString(userecRaw.LastHost[:])

	return user
}
