package initgin

import (
	"runtime"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

var (
	testUserec          *bbs.Userec
	testUserecRaw3      *ptttype.UserecRaw
	testNewPostUserRaw1 *ptttype.UserecRaw
)

func initTestVars() {
	if testUserec != nil {
		return
	}

	testUserec = &bbs.Userec{
		Version:  4194,
		UUserID:  bbs.UUserID("SYSOP"),
		Username: "SYSOP",
		Realname: []byte{ // CodingMan
			0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e,
		},
		Nickname: []byte{0xaf, 0xab}, // 神

		Uflag:        33557088,
		Userlevel:    536871943,
		Numlogindays: 2,
		Numposts:     0,
		Firstlogin:   1600681288,
		Lastlogin:    1600756094,
		Lasthost:     "59.124.167.226",
		/*
			Address: []byte{ //新竹縣子虛鄉烏有村543號
				0xb7, 0x73, 0xa6, 0xcb, 0xbf, 0xa4, 0xa4, 0x6c, 0xb5, 0xea,
				0xb6, 0x6d, 0xaf, 0x51, 0xa6, 0xb3, 0xa7, 0xf8, 0x35, 0x34,
				0x33, 0xb8, 0xb9,
			},
		*/
		Over18:   true,
		Pager:    ptttype.PAGER_ON,
		Career:   []byte{0xa5, 0xfe, 0xb4, 0xba, 0xb3, 0x6e, 0xc5, 0xe9}, // 全景軟體
		LastSeen: 1600681288,
	}

	testUserecRaw3 = &ptttype.UserecRaw{
		Version: ptttype.PASSWD_VERSION,
		UserID:  ptttype.UserID_t{'t', 'e', 's', 't'}, // test
		RealName: ptttype.RealName_t{ // CodingMan
			0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e,
		},
		Nickname:   ptttype.Nickname_t{0xaf, 0xab}, // 神
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    7 | ptttype.PERM_BOARD | ptttype.PERM_POST | ptttype.PERM_LOGINOK,
		NumLoginDays: 2,
		NumPosts:     0,
		FirstLogin:   1600681288,
		LastLogin:    1600756094,
		LastHost: ptttype.IPv4_t{ // 59.124.167.226
			0x35, 0x39, 0x2e, 0x31, 0x32, 0x34, 0x2e, 0x31, 0x36, 0x37,
			0x2e, 0x32, 0x32, 0x36,
		},
		Address: ptttype.Address_t{ // 新竹縣子虛鄉烏有村543號
			0xb7, 0x73, 0xa6, 0xcb, 0xbf, 0xa4, 0xa4, 0x6c, 0xb5, 0xea,
			0xb6, 0x6d, 0xaf, 0x51, 0xa6, 0xb3, 0xa7, 0xf8, 0x35, 0x34,
			0x33, 0xb8, 0xb9,
		},
		Over18:   true,
		Pager:    ptttype.PAGER_ON,
		Career:   ptttype.Career_t{0xa5, 0xfe, 0xb4, 0xba, 0xb3, 0x6e, 0xc5, 0xe9}, // 全景軟體
		LastSeen: 1600681288,
	}

	testNewPostUserRaw1 = &ptttype.UserecRaw{
		Version:    4194,
		UserID:     ptttype.UserID_t{65, 49}, // A1
		RealName:   ptttype.RealName_t{67, 111, 100, 105, 110, 103, 77, 97, 110},
		Nickname:   ptttype.Nickname_t{175, 171},
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    7 | ptttype.PERM_LOGINOK | ptttype.PERM_POST,
		NumLoginDays: 2,
		NumPosts:     0,
		FirstLogin:   1600681288,
		LastLogin:    1600756094,
		LastHost:     ptttype.IPv4_t{53, 57, 46, 49, 50, 52, 46, 49, 54, 55, 46, 50, 50, 54},
		Address:      ptttype.Address_t{183, 115, 166, 203, 191, 164, 164, 108, 181, 234, 182, 109, 175, 81, 166, 179, 167, 248, 53, 52, 51, 184, 185},
		Over18:       true,
		Pager:        ptttype.PAGER_ON,
		Career:       ptttype.Career_t{165, 254, 180, 186, 179, 110, 197, 233},
		LastSeen:     1600681288,
	}
}

func freeTestVars() {
	testUserec = nil
	testUserecRaw3 = nil
	testNewPostUserRaw1 = nil

	runtime.GC()
}
