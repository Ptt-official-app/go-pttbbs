package bbs

import (
	"runtime"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/sirupsen/logrus"
)

var (
	testUserecRaw   *ptttype.UserecRaw
	testUserecRaw3  *ptttype.UserecRaw
	testUserec1     *Userec
	testUserec2     *Userec
	testUserec3     *Userec
	testUserec4     *Userec
	testUserec5     *Userec
	testUserecEmpty *Userec
	testUserec6     *Userec

	testNewPostUserRaw1 *ptttype.UserecRaw

	testOpenUserecFile1     []*Userec = nil
	TEST_N_OPEN_USER_FILE_1           = 50

	testBoardSummaryRaw6 *ptttype.BoardSummaryRaw

	testBoardSummary1  *BoardSummary
	testBoardSummary6  *BoardSummary
	testBoardSummary7  *BoardSummary
	testBoardSummary8  *BoardSummary
	testBoardSummary9  *BoardSummary
	testBoardSummary10 *BoardSummary
	testBoardSummary11 *BoardSummary
	testBoardSummary13 *BoardSummary

	testBoardDetail3  *BoardDetail
	testBoardDetail1  *BoardDetail
	testClassDetail2  *BoardDetail
	testClassDetail5  *BoardDetail
	testBoardDetail12 *BoardDetail
	testBoardDetail6  *BoardDetail
	testBoardDetail7  *BoardDetail

	testBoardDetailRaw3 *ptttype.BoardDetailRaw
	testBoardHeader3    *ptttype.BoardHeaderRaw

	testClassSummary2 *BoardSummary
	testClassSummary5 *BoardSummary

	testArticleSummary0 *ArticleSummary
	testArticleSummary1 *ArticleSummary

	testBottomSummary1 *ArticleSummary

	testContent1 []byte

	testPermissionUserecRaw  *ptttype.UserecRaw
	testPermissionUserecRaw2 *ptttype.UserecRaw
	testPermissionUserecRaw3 *ptttype.UserecRaw
)

func initTestVars() {
	if testUserecRaw != nil {
		if testBoardSummary6.BoardClass == nil || testBoardSummary6.BoardClass[0] == 0 {
			logrus.Errorf("initVars: invalid testBoardSummary6: %v", testBoardSummary6)
		}
		return
	}

	testUserecRaw = &ptttype.UserecRaw{
		Version: ptttype.PASSWD_VERSION,
		UserID:  ptttype.UserID_t{0x53, 0x59, 0x53, 0x4f, 0x50}, // SYSOP
		RealName: ptttype.RealName_t{ // CodingMan
			0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e,
		},
		Nickname:   ptttype.Nickname_t{0xaf, 0xab}, // 神
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    536871943,
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

	testPermissionUserecRaw = &ptttype.UserecRaw{
		Version: ptttype.PASSWD_VERSION,
		UserID:  ptttype.UserID_t{'t', 'e', 's', 't'}, // test
		RealName: ptttype.RealName_t{ // CodingMan
			0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e,
		},
		Nickname:   ptttype.Nickname_t{0xaf, 0xab}, // 神
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    ptttype.PERM_SYSOP,
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

	testPermissionUserecRaw2 = &ptttype.UserecRaw{
		Version: ptttype.PASSWD_VERSION,
		UserID:  ptttype.UserID_t{'t', 'e', 's', 't', '1'}, // test
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

	testPermissionUserecRaw3 = &ptttype.UserecRaw{
		Version: ptttype.PASSWD_VERSION,
		UserID:  ptttype.UserID_t{'t', 'e', 's', 't', '2'}, // test
		RealName: ptttype.RealName_t{ // CodingMan
			0x43, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x4d, 0x61, 0x6e,
		},
		Nickname:   ptttype.Nickname_t{0xaf, 0xab}, // 神
		PasswdHash: ptttype.Passwd_t{0x62, 0x68, 0x77, 0x76, 0x4f, 0x4a, 0x74, 0x66, 0x54, 0x31, 0x54, 0x41, 0x49, 0x00},

		UFlag:        33557088,
		UserLevel:    ptttype.PERM_BASIC,
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
		PasswdHash: ptttype.Passwd_t{98, 104, 119, 118, 79, 74, 116, 102, 84, 49, 84, 65, 73, 0},

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

	testUserec1 = &Userec{
		Version:  4194,
		UUserID:  UUserID("SYSOP"),
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

	testUserec6 = &Userec{
		Version:      ptttype.PASSWD_VERSION,
		UUserID:      UUserID("B1"),
		Username:     "B1",
		Lasthost:     "127.0.0.1",
		Uflag:        33557088,
		Userlevel:    7,
		Numlogindays: 1,
		Pager:        1,
		Over18:       true,
	}

	testUserecEmpty = &Userec{}

	testBoardSummaryRaw6 = &ptttype.BoardSummaryRaw{
		Gid:     5,
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

	testBoardSummary1 = &BoardSummary{
		Gid:      2,
		Bid:      1,
		BBoardID: BBoardID("1_SYSOP"),
		BrdAttr:  ptttype.BRD_POSTMASK,
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "SYSOP",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xaf, 0xb8, 0xaa, 0xf8, 0xa6, 0x6e, 0x21,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "SYSOP",
		IdxByClass: "vFSt-Q@SYSOP",
	}

	testBoardDetail1 = &BoardDetail{
		Gid:      2,
		Bid:      1,
		BBoardID: BBoardID("1_SYSOP"),
		BrdAttr:  ptttype.BRD_POSTMASK,
		Brdname:  "SYSOP",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xaf, 0xb8, 0xaa, 0xf8, 0xa6, 0x6e, 0x21,
		},
		BoardType:        []byte{0xa1, 0xb7},
		BM:               []UUserID{},
		IdxByName:        "SYSOP",
		IdxByClass:       "vFSt-Q@SYSOP",
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
	}

	testBoardSummary6 = &BoardSummary{
		Gid:      5,
		Bid:      6,
		BBoardID: BBoardID("6_ALLPOST"),
		BrdAttr:  ptttype.BRD_POSTMASK,
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "ALLPOST",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xb8, 0xf3, 0xaa, 0x4f, 0xa6, 0xa1, 0x4c, 0x4f, 0x43, 0x41,
			0x4c, 0xb7, 0x73, 0xa4, 0xe5, 0xb3, 0xb9,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "ALLPOST",
		IdxByClass: "vFSt-Q@ALLPOST",
	}

	testBoardDetail6 = &BoardDetail{
		Gid:      5,
		Bid:      6,
		BBoardID: BBoardID("6_ALLPOST"),
		BrdAttr:  ptttype.BRD_POSTMASK,
		Brdname:  "ALLPOST",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xb8, 0xf3, 0xaa, 0x4f, 0xa6, 0xa1, 0x4c, 0x4f, 0x43, 0x41,
			0x4c, 0xb7, 0x73, 0xa4, 0xe5, 0xb3, 0xb9,
		},
		BoardType:        []byte{0xa1, 0xb7},
		BM:               []UUserID{},
		IdxByName:        "ALLPOST",
		IdxByClass:       "vFSt-Q@ALLPOST",
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
		Level:            ptttype.PERM_SYSOP,
	}

	testBoardSummary7 = &BoardSummary{
		Gid:      5,
		Bid:      7,
		BBoardID: BBoardID("7_deleted"),
		Brdname:  "deleted",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		StatAttr: ptttype.NBRD_FAV,
		RealTitle: []byte{
			0xb8, 0xea, 0xb7, 0xbd, 0xa6, 0x5e, 0xa6, 0xac, 0xb5, 0xa9,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "deleted",
		IdxByClass: "vFSt-Q@deleted",
	}

	testBoardDetail7 = &BoardDetail{
		Gid:      5,
		Bid:      7,
		BBoardID: BBoardID("7_deleted"),
		Brdname:  "deleted",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xb8, 0xea, 0xb7, 0xbd, 0xa6, 0x5e, 0xa6, 0xac, 0xb5, 0xa9,
		},
		BoardType:        []byte{0xa1, 0xb7},
		BM:               []UUserID{},
		IdxByName:        "deleted",
		IdxByClass:       "vFSt-Q@deleted",
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
	}

	testBoardSummary8 = &BoardSummary{
		Gid:      5,
		Bid:      8,
		BBoardID: BBoardID("8_Note"),
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "Note",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xb0, 0xca, 0xba, 0x41, 0xac, 0xdd, 0xaa, 0x4f, 0xa4, 0xce,
			0xba, 0x71, 0xa6, 0xb1, 0xa7, 0xeb, 0xbd, 0x5a,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "Note",
		IdxByClass: "vFSt-Q@Note",
	}

	testBoardSummary9 = &BoardSummary{
		Gid:      5,
		Bid:      9,
		BBoardID: BBoardID("9_Record"),
		BrdAttr:  ptttype.BRD_POSTMASK,
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "Record",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xa7, 0xda, 0xad, 0xcc, 0xaa, 0xba, 0xa6, 0xa8, 0xaa, 0x47,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "Record",
		IdxByClass: "vFSt-Q@Record",
	}

	testBoardSummary10 = &BoardSummary{
		Gid:      5,
		Bid:      10,
		BBoardID: BBoardID("10_WhoAmI"),
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "WhoAmI",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xa8, 0xfe, 0xa8, 0xfe, 0xa1, 0x41, 0xb2, 0x71, 0xb2, 0x71,
			0xa7, 0xda, 0xac, 0x4f, 0xbd, 0xd6, 0xa1, 0x49,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "WhoAmI",
		IdxByClass: "vFSt-Q@WhoAmI",
	}

	testBoardSummary11 = &BoardSummary{
		Gid:      5,
		Bid:      11,
		BBoardID: BBoardID("11_EditExp"),
		StatAttr: ptttype.NBRD_FAV,
		Brdname:  "EditExp",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xbd, 0x64, 0xa5, 0xbb, 0xba, 0xeb, 0xc6, 0x46, 0xa7, 0xeb,
			0xbd, 0x5a, 0xb0, 0xcf,
		},
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "EditExp",
		IdxByClass: "vFSt-Q@EditExp",
	}

	testBoardDetail12 = &BoardDetail{
		Gid:      5,
		Bid:      12,
		BBoardID: BBoardID("12_ALLHIDPOST"),
		BrdAttr:  ptttype.BRD_POSTMASK | ptttype.BRD_HIDE,
		Brdname:  "ALLHIDPOST",
		BoardClass: []byte{
			0xbc, 0x54, 0xad, 0xf9,
		},
		RealTitle: []byte{
			0xb8, 0xf3, 0xaa, 0x4f, 0xa6, 0xa1, 0x4c, 0x4f, 0x43, 0x41,
			0x4c, 0xb7, 0x73, 0xa4, 0xe5, 0xb3, 0xb9, 0x28, 0xc1, 0xf4,
			0xaa, 0x4f, 0x29,
		},
		BoardType:        []byte{0xa1, 0xb7},
		BM:               []UUserID{},
		IdxByName:        "ALLHIDPOST",
		IdxByClass:       "vFSt-Q@ALLHIDPOST",
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
		Level:            ptttype.PERM_SYSOP,
	}

	testBoardSummary13 = &BoardSummary{
		Gid:        5,
		Bid:        13,
		BBoardID:   BBoardID("13_mnewboard"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "mnewboard",
		BoardClass: []byte("CPBL"),
		RealTitle:  []byte("new-board"),
		BoardType:  []byte{0xa1, 0xb7},
		BM:         []UUserID{},
		IdxByName:  "mnewboard",
		IdxByClass: "Q1BCTA@mnewboard",
		BrdAttr:    0x200000,
	}

	testClassSummary2 = &BoardSummary{
		Gid:        1,
		Bid:        2,
		BBoardID:   BBoardID("2_1..........."),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "1...........",
		BoardClass: []byte("...."),
		RealTitle: []byte{
			0xa4, 0xa4, 0xa5, 0xa1, 0xac, 0x46, 0xa9, 0xb2,
			0x20, 0x20, 0xa1, 0x6d, 0xb0, 0xaa, 0xc0, 0xa3, 0xa6, 0x4d,
			0xc0, 0x49, 0x2c, 0xab, 0x44, 0xa4, 0x48, 0xa5, 0x69, 0xbc,
			0xc4, 0xa1, 0x6e,
		},
		BoardType:  []byte{0xa3, 0x55},
		BM:         []UUserID{},
		IdxByName:  "1...........",
		IdxByClass: "Li4uLg@1...........",
		BrdAttr:    0x000008,
	}

	testClassDetail2 = &BoardDetail{
		Gid:        1,
		Bid:        2,
		BBoardID:   BBoardID("2_1..........."),
		Brdname:    "1...........",
		BoardClass: []byte("...."),
		RealTitle: []byte{
			0xa4, 0xa4, 0xa5, 0xa1, 0xac, 0x46, 0xa9, 0xb2,
			0x20, 0x20, 0xa1, 0x6d, 0xb0, 0xaa, 0xc0, 0xa3, 0xa6, 0x4d,
			0xc0, 0x49, 0x2c, 0xab, 0x44, 0xa4, 0x48, 0xa5, 0x69, 0xbc,
			0xc4, 0xa1, 0x6e,
		},
		BoardType:        []byte{0xa3, 0x55},
		BM:               []UUserID{},
		IdxByName:        "1...........",
		IdxByClass:       "Li4uLg@1...........",
		BrdAttr:          ptttype.BRD_GROUPBOARD,
		Level:            ptttype.PERM_SYSOP,
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
	}

	testClassSummary5 = &BoardSummary{
		Gid:        1,
		Bid:        5,
		BBoardID:   BBoardID("5_2..........."),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "2...........",
		BoardClass: []byte("...."),
		RealTitle: []byte{
			0xa5, 0xab, 0xa5, 0xc1, 0xbc, 0x73, 0xb3, 0xf5,
			0x20, 0x20, 0x20, 0x20, 0x20, 0xb3, 0xf8, 0xa7, 0x69, 0x20,
			0x20, 0xaf, 0xb8, 0xaa, 0xf8, 0x20, 0x20, 0xa3, 0xad, 0xa1,
			0x49,
		},
		BoardType:  []byte{0xa3, 0x55},
		BM:         []UUserID{},
		IdxByName:  "2...........",
		IdxByClass: "Li4uLg@2...........",
		BrdAttr:    0x000008,
	}

	testClassDetail5 = &BoardDetail{
		Gid:        1,
		Bid:        5,
		BBoardID:   BBoardID("5_2..........."),
		Brdname:    "2...........",
		BoardClass: []byte("...."),
		RealTitle: []byte{
			0xa5, 0xab, 0xa5, 0xc1, 0xbc, 0x73, 0xb3, 0xf5,
			0x20, 0x20, 0x20, 0x20, 0x20, 0xb3, 0xf8, 0xa7, 0x69, 0x20,
			0x20, 0xaf, 0xb8, 0xaa, 0xf8, 0x20, 0x20, 0xa3, 0xad, 0xa1,
			0x49,
		},
		BoardType:        []byte{0xa3, 0x55},
		BM:               []UUserID{},
		IdxByName:        "2...........",
		IdxByClass:       "Li4uLg@2...........",
		BrdAttr:          ptttype.BRD_GROUPBOARD,
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
	}

	testBoardHeader3 = &ptttype.BoardHeaderRaw{
		Brdname: ptttype.BoardID_t{'S', 'Y', 'S', 'O', 'P'},
		Title: ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xaf, 0xb8, 0xaa,
			0xf8, 0xa6, 0x6e, 0x21,
		},
		BrdAttr: ptttype.BRD_POSTMASK,
		Gid:     2,
	}
	testBoardDetailRaw3 = &ptttype.BoardDetailRaw{
		Bid:            1,
		BoardHeaderRaw: testBoardHeader3,
	}

	testBoardDetail3 = &BoardDetail{
		Brdname: "SYSOP",
		RealTitle: []byte{
			0xaf, 0xb8, 0xaa, 0xf8, 0xa6, 0x6e, 0x21,
		},
		BoardClass:       []byte{0xbc, 0x54, 0xad, 0xf9},
		BoardType:        []byte{0xa1, 0xb7},
		BM:               []UUserID{},
		BrdAttr:          ptttype.BRD_POSTMASK,
		Gid:              2,
		Bid:              1,
		BBoardID:         "1_SYSOP",
		IdxByName:        "SYSOP",
		IdxByClass:       "vFSt-Q@SYSOP",
		PostType:         [][]byte{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		PostTypeTemplate: []bool{false, false, false, false, false, false, false, false},
	}

	testArticleSummary0 = &ArticleSummary{
		BBoardID:   BBoardID("10_WhoAmI"),
		ArticleID:  "1Vo_M_CD",
		IsDeleted:  false,
		Filename:   "M.1607202239.A.30D",
		CreateTime: 1607202239,
		MTime:      1607202238,
		Owner:      "SYSOP",
		FullTitle: []byte{
			0x5b, 0xb0, 0xdd, 0xc3, 0x44, 0x5d, 0x20, 0xa7,
			0xda, 0xac, 0x4f, 0xbd, 0xd6, 0xa1, 0x48, 0xa1,
			0xe3,
		},

		Class:     []byte{0xb0, 0xdd, 0xc3, 0x44},
		Idx:       "1607202239@1Vo_M_CD",
		RealTitle: []byte{0xa7, 0xda, 0xac, 0x4f, 0xbd, 0xd6, 0xa1, 0x48, 0xa1, 0xe3},
	}

	testArticleSummary1 = &ArticleSummary{
		BBoardID:   BBoardID("10_WhoAmI"),
		ArticleID:  "1Vo_f30D",
		IsDeleted:  false,
		Filename:   "M.1607203395.A.00D",
		CreateTime: 1607203395,
		MTime:      1607203394,
		Owner:      "SYSOP",
		FullTitle: []byte{
			0x5b, 0xa4, 0xdf, 0xb1, 0x6f, 0x5d, 0x20, 0xb5,
			0x4d, 0xab, 0xe1, 0xa9, 0x4f, 0xa1, 0x48, 0xa1,
			0xe3,
		},

		Filemode: ptttype.FILE_MARKED,

		Class: []byte{0xa4, 0xdf, 0xb1, 0x6f},
		Idx:   "1607203395@1Vo_f30D",
		RealTitle: []byte{
			0xb5, 0x4d, 0xab, 0xe1, 0xa9, 0x4f, 0xa1, 0x48, 0xa1, 0xe3,
		},
	}

	testBottomSummary1 = &ArticleSummary{
		BBoardID:   BBoardID("10_WhoAmI"),
		ArticleID:  "1Vo_f30D",
		IsDeleted:  false,
		Filename:   "M.1607203395.A.00D",
		CreateTime: 1607203395,
		MTime:      1607203394,
		Owner:      "SYSOP",
		FullTitle: []byte{
			0x5b, 0xa4, 0xdf, 0xb1, 0x6f, 0x5d, 0x20, 0xb5,
			0x4d, 0xab, 0xe1, 0xa9, 0x4f, 0xa1, 0x48, 0xa1,
			0xe3,
		},

		Filemode: ptttype.FILE_MULTI,
		Money:    -2147483646,

		Class: []byte{0xa4, 0xdf, 0xb1, 0x6f},
		Idx:   "1607203395@1Vo_f30D",
		RealTitle: []byte{
			0xb5, 0x4d, 0xab, 0xe1, 0xa9, 0x4f, 0xa1, 0x48, 0xa1, 0xe3,
		},
	}

	testContent1 = []byte{
		0xa7, 0x40, 0xaa, 0xcc, 0x3a, 0x20, 0x53, 0x59, 0x53,
		0x4f, 0x50, 0x20, 0x28, 0x29, 0x20, 0xac, 0xdd, 0xaa,
		0x4f, 0x3a, 0x20, 0x57, 0x68, 0x6f, 0x41, 0x6d, 0x49,
		0x0a, 0xbc, 0xd0, 0xc3, 0x44, 0x3a, 0x20, 0x5b, 0xb0,
		0xdd, 0xc3, 0x44, 0x5d, 0x20, 0xa7, 0xda, 0xac, 0x4f,
		0xbd, 0xd6, 0xa1, 0x48, 0xa1, 0xe3, 0x0a, 0xae, 0xc9,
		0xb6, 0xa1, 0x3a, 0x20, 0x53, 0x75, 0x6e, 0x20, 0x44,
		0x65, 0x63, 0x20, 0x20, 0x36, 0x20, 0x30, 0x35, 0x3a,
		0x30, 0x33, 0x3a, 0x35, 0x37, 0x20, 0x32, 0x30, 0x32,
		0x30, 0x0a, 0x0a, 0xa7, 0xda, 0xac, 0x4f, 0xbd, 0xd6,
		0xa1, 0x48, 0xa1, 0xe3, 0x0a, 0x0a, 0xa7, 0xda, 0xa6,
		0x62, 0xad, 0xfe, 0xb8, 0xcc, 0xa1, 0x48, 0xa1, 0xe3,
		0x0a, 0x0a, 0xa7, 0xda, 0xac, 0xb0, 0xa4, 0xb0, 0xbb,
		0xf2, 0xb7, 0x7c, 0xa6, 0x62, 0xb3, 0x6f, 0xb8, 0xcc,
		0xa9, 0x4f, 0xa1, 0x48, 0xa1, 0xe3, 0x0a, 0x0a, 0x2d,
		0x2d, 0x0a, 0xa1, 0xb0, 0x20, 0xb5, 0x6f, 0xab, 0x48,
		0xaf, 0xb8, 0x3a, 0x20, 0xa7, 0xe5, 0xbd, 0xf0, 0xbd,
		0xf0, 0x20, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72, 0x28,
		0x70, 0x74, 0x74, 0x64, 0x6f, 0x63, 0x6b, 0x65, 0x72,
		0x2e, 0x74, 0x65, 0x73, 0x74, 0x29, 0x2c, 0x20, 0xa8,
		0xd3, 0xa6, 0xdb, 0x3a, 0x20, 0x31, 0x37, 0x32, 0x2e,
		0x31, 0x38, 0x2e, 0x30, 0x2e, 0x31, 0x0a,
	}

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

func freeTestVars() {
	testUserecRaw = nil
	testUserecRaw3 = nil
	testUserec1 = nil
	testUserec2 = nil
	testUserec3 = nil
	testUserec4 = nil
	testUserec5 = nil
	testUserecEmpty = nil
	testUserec6 = nil

	testNewPostUserRaw1 = nil

	testOpenUserecFile1 = nil

	testBoardSummaryRaw6 = nil

	testBoardSummary1 = nil
	testBoardSummary6 = nil
	testBoardSummary7 = nil
	testBoardSummary8 = nil
	testBoardSummary9 = nil
	testBoardSummary10 = nil
	testBoardSummary11 = nil
	testBoardSummary13 = nil
	testClassSummary2 = nil
	testClassSummary5 = nil

	testArticleSummary0 = nil
	testArticleSummary1 = nil

	testBottomSummary1 = nil

	testContent1 = nil

	testPermissionUserecRaw = nil
	testPermissionUserecRaw2 = nil
	testPermissionUserecRaw3 = nil
	runtime.GC()
}
