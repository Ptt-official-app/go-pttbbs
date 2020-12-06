package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

var (
	testUserecRaw = &ptttype.UserecRaw{
		Version:    ptttype.PASSWD_VERSION,
		UserID:     ptttype.UserID_t{0x53, 0x59, 0x53, 0x4f, 0x50},                           // SYSOP
		RealName:   ptttype.RealName_t{0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e}, // CodingMan
		Nickname:   ptttype.Nickname_t{0xaf, 0xab},                                           // 神
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    536871943,
		NumLoginDays: 2,
		NumPosts:     0,
		FirstLogin:   1600681288,
		LastLogin:    1600756094,
		LastHost:     ptttype.IPv4_t{0x35, 0x39, 0x2e, 0x31, 0x32, 0x34, 0x2e, 0x31, 0x36, 0x37, 0x2e, 0x32, 0x32, 0x36},                                                          //59.124.167.226
		Address:      ptttype.Address_t{0xb7, 0x73, 0xa6, 0xcb, 0xbf, 0xa4, 0xa4, 0x6c, 0xb5, 0xea, 0xb6, 0x6d, 0xaf, 0x51, 0xa6, 0xb3, 0xa7, 0xf8, 0x35, 0x34, 0x31, 0xb8, 0xb9}, //新竹縣子虛鄉烏有村543號
		Over18:       true,
		Pager:        ptttype.PAGER_ON,
		Career:       ptttype.Career_t{0xa5, 0xfe, 0xb4, 0xba, 0xb3, 0x6e, 0xc5, 0xe9}, //全景軟體
		LastSeen:     1600681288,
	}

	testUserec1 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "SYSOP",
		Realname:     "CodingMan",
		Nickname:     "神",
		Passwd:       "bhwvOJtfT1TAI",
		Uflag:        33557088,
		Userlevel:    536871943,
		Numlogindays: 2,
		Numposts:     0,
		Firstlogin:   1600681288,
		Lastlogin:    1600756094,
		Lasthost:     "59.124.167.226",
	}

	testUserec2 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "CodingMan",
		Realname:     "朱元璋",
		Nickname:     "程式俠",
		Passwd:       "u8mLG.ktfOk3w",
		Uflag:        33557216,
		Userlevel:    31,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600737659,
		Lastlogin:    1600737960,
		Lasthost:     "59.124.167.226",
	}

	testUserec3 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "pichu",
		Realname:     "Pichu",
		Nickname:     "Pichu",
		Passwd:       "KO27TyME.3/tw",
		Uflag:        33557216,
		Userlevel:    7,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600755675,
		Lastlogin:    1600766204,
		Lasthost:     "103.246.218.43",
	}

	testUserec4 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "Kahou",
		Realname:     "林嘉豪",
		Nickname:     "Kahou",
		Passwd:       "V3nkaYTLnDPUA",
		Uflag:        33557216,
		Userlevel:    7,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600758266,
		Lastlogin:    1600758266,
		Lasthost:     "180.217.174.18",
	}

	testUserec5 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "Kahou2",
		Realname:     "Kahou",
		Nickname:     "kahou",
		Passwd:       "R7shIAOZgQCKs",
		Uflag:        33557216,
		Userlevel:    31,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600758939,
		Lastlogin:    1600760401,
		Lasthost:     "180.217.174.18",
	}
	testUserecEmpty = &Userec{}

	testUserec6 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "B1",
		Lasthost:     "127.0.0.1",
		Uflag:        33557088,
		Userlevel:    7,
		Numlogindays: 1,
	}

	testOpenUserecFile1     []*Userec = nil
	TEST_N_OPEN_USER_FILE_1           = 50
)

func initTestVars() {
	if testOpenUserecFile1 == nil {
		testOpenUserecFile1 = make([]*Userec, TEST_N_OPEN_USER_FILE_1)
		for i := 0; i < TEST_N_OPEN_USER_FILE_1; i++ {
			testOpenUserecFile1[i] = testUserecEmpty
		}
		testOpenUserecFile1[0] = testUserec1
		testOpenUserecFile1[1] = testUserec2
		testOpenUserecFile1[2] = testUserec3
		testOpenUserecFile1[3] = testUserec4
		testOpenUserecFile1[4] = testUserec5
	}
}

func freeTestVars() {
	testOpenUserecFile1 = nil
}
