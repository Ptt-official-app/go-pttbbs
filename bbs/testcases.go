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
		Address:      ptttype.Address_t{0xb7, 0x73, 0xa6, 0xcb, 0xbf, 0xa4, 0xa4, 0x6c, 0xb5, 0xea, 0xb6, 0x6d, 0xaf, 0x51, 0xa6, 0xb3, 0xa7, 0xf8, 0x35, 0x34, 0x33, 0xb8, 0xb9}, //新竹縣子虛鄉烏有村543號
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
		Uflag:        33557088,
		Userlevel:    536871943,
		Numlogindays: 2,
		Numposts:     0,
		Firstlogin:   1600681288,
		Lastlogin:    1600756094,
		Lasthost:     "59.124.167.226",
		Address:      "新竹縣子虛鄉烏有村543號",
		Over18:       true,
		Pager:        1,
		Career:       "全景軟體",
	}

	testUserec2 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "CodingMan",
		Realname:     "朱元璋",
		Nickname:     "程式俠",
		Uflag:        33557216,
		Userlevel:    31,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600737659,
		Lastlogin:    1600737960,
		Lasthost:     "59.124.167.226",
		Email:        "x",
		Address:      "新竹縣子虛鄉烏有村543號",
		Justify:      "[SYSOP] 09/22/2020 01:25:53 Tue",
		Over18:       true,
		Pager:        1,
		Career:       "全景軟體",
	}

	testUserec3 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "pichu",
		Realname:     "Pichu",
		Nickname:     "Pichu",
		Uflag:        33557216,
		Userlevel:    7,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600755675,
		Lastlogin:    1600766204,
		Lasthost:     "103.246.218.43",
		Email:        "pichu@tih.tw",
		Address:      "北市蘆洲區123路3號",
		Justify:      "<Email>",
		Over18:       true,
		Pager:        1,
		Career:       "台灣智慧家庭",
	}

	testUserec4 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "Kahou",
		Realname:     "林嘉豪",
		Nickname:     "Kahou",
		Uflag:        33557216,
		Userlevel:    7,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600758266,
		Lastlogin:    1600758266,
		Lasthost:     "180.217.174.18",
		Email:        "creator.kahou@gmail.com",
		Address:      "新北市板橋信義路111號",
		Justify:      "<Email>",
		Over18:       true,
		Pager:        1,
		Career:       "我的服務單位",
	}

	testUserec5 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "Kahou2",
		Realname:     "Kahou",
		Nickname:     "kahou",
		Uflag:        33557216,
		Userlevel:    31,
		Numlogindays: 1,
		Numposts:     0,
		Firstlogin:   1600758939,
		Lastlogin:    1600760401,
		Lasthost:     "180.217.174.18",
		Email:        "x",
		Address:      "新北市板橋區信義路111號",
		Justify:      "[SYSOP] 09/22/2020 07:51:12 Tue",
		Over18:       true,
		Pager:        1,
		Career:       "我的服務單位",
	}
	testUserecEmpty = &Userec{}

	testUserec6 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		Userid:       "B1",
		Lasthost:     "127.0.0.1",
		Uflag:        33557088,
		Userlevel:    7,
		Numlogindays: 1,
		Pager:        1,
		Over18:       true,
	}

	testOpenUserecFile1     []*Userec = nil
	TEST_N_OPEN_USER_FILE_1           = 50

	testBoardSummaryRaw6 = &ptttype.BoardSummaryRaw{
		Bid:     6,
		BrdAttr: ptttype.BRD_POSTMASK,
		Brdname: &ptttype.BoardID_t{'A', 'L', 'L', 'P', 'O', 'S', 'T', 0x00, 0x2e, 0x2e, 0x2e, 0x2e},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xb8, 0xf3, 0xaa,
			0x4f, 0xa6, 0xa1, 0x4c, 0x4f, 0x43, 0x41, 0x4c, 0xb7, 0x73,
			0xa4, 0xe5, 0xb3, 0xb9, 0x00, 0x20, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x00, 0x6e,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardSummary6 = &BoardSummary{
		BBoardID:   BBoardID("6_ALLPOST"),
		BrdAttr:    ptttype.BRD_POSTMASK,
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "ALLPOST",
		BoardClass: "嘰哩",
		RealTitle:  "跨板式LOCAL新文章",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary7 = &BoardSummary{
		BBoardID:   BBoardID("7_deleted"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "deleted",
		BoardClass: "嘰哩",
		RealTitle:  "資源回收筒",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary11 = &BoardSummary{
		BBoardID:   BBoardID("11_EditExp"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "EditExp",
		BoardClass: "嘰哩",
		RealTitle:  "範本精靈投稿區",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary8 = &BoardSummary{
		BBoardID:   BBoardID("8_Note"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "Note",
		BoardClass: "嘰哩",
		RealTitle:  "動態看板及歌曲投稿",
		BoardType:  "◎",
		BM:         []string{},
	}

	testArticleSummary0 = &ArticleSummary{
		BBoardID:   BBoardID("10_WhoAmI"),
		ArticleID:  "1_1Vo_M_CD",
		IsDeleted:  false,
		Filename:   "M.1607202239.A.30D",
		CreateTime: 1607202239,
		MTime:      1607202238,
		Owner:      "SYSOP",
		Date:       "12/06",
		Title:      "[問題] 我是誰？～",
		URL:        "http://localhost/bbs/WhoAmI/M.1607202239.A.30D.html",
	}

	testArticleSummary1 = &ArticleSummary{
		BBoardID:   BBoardID("10_WhoAmI"),
		ArticleID:  "2_1Vo_f30D",
		IsDeleted:  false,
		Filename:   "M.1607203395.A.00D",
		CreateTime: 1607203395,
		MTime:      1607203394,
		Owner:      "SYSOP",
		Date:       "12/06",
		Title:      "[心得] 然後呢？～",
		Filemode:   ptttype.FILE_MARKED,
		URL:        "http://localhost/bbs/WhoAmI/M.1607203395.A.00D.html",
	}
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
}
