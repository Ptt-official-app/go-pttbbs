package ptt

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

var (
	testUserecRaw1 = &ptttype.UserecRaw{
		Version:    4194,
		UserID:     ptttype.UserID_t{83, 89, 83, 79, 80},
		RealName:   ptttype.RealName_t{67, 111, 100, 105, 110, 103, 77, 97, 110},
		Nickname:   ptttype.Nickname_t{175, 171},
		PasswdHash: ptttype.Passwd_t{98, 104, 119, 118, 79, 74, 116, 102, 84, 49, 84, 65, 73, 0},

		UFlag:        33557088,
		UserLevel:    536871943,
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

	testUserecRaw2 = &ptttype.UserecRaw{
		Version:      ptttype.PASSWD_VERSION,
		UserID:       ptttype.UserID_t{'C', 'o', 'd', 'i', 'n', 'g', 'M', 'a', 'n'},
		UFlag:        33557216,
		UserLevel:    7,
		NumLoginDays: 1,
		NumPosts:     0,
		FirstLogin:   1600737659,
		LastLogin:    1600737960,
		LastHost:     ptttype.IPv4_t{'5', '9', '.', '1', '2', '4', '.', '1', '6', '7', '.', '2', '2', '6'},
	}

	testSetupNewUser1 = &ptttype.UserecRaw{
		Version:    4194,
		UserID:     ptttype.UserID_t{65, 48}, //A0
		RealName:   ptttype.RealName_t{67, 111, 100, 105, 110, 103, 77, 97, 110},
		Nickname:   ptttype.Nickname_t{175, 171},
		PasswdHash: ptttype.Passwd_t{98, 104, 119, 118, 79, 74, 116, 102, 84, 49, 84, 65, 73, 0},

		UFlag:        33557088,
		UserLevel:    536871943,
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

	testNewRegister1 = &ptttype.UserecRaw{
		Version:    4194,
		UserID:     ptttype.UserID_t{66, 49}, //B1
		RealName:   ptttype.RealName_t{67, 111, 100, 105, 110, 103, 77, 97, 110},
		Nickname:   ptttype.Nickname_t{175, 171},
		PasswdHash: ptttype.Passwd_t{98, 104, 119, 118, 79, 74, 116, 102, 84, 49, 84, 65, 73, 0},

		UFlag:        33557088,
		UserLevel:    7,
		NumLoginDays: 1,
		NumPosts:     0,
		FirstLogin:   1600681288,
		LastLogin:    1600756094,
		LastHost:     ptttype.IPv4_t{53, 57, 46, 49, 50, 52, 46, 49, 54, 55, 46, 50, 50, 54},
		Address:      ptttype.Address_t{183, 115, 166, 203, 191, 164, 164, 108, 181, 234, 182, 109, 175, 81, 166, 179, 167, 248, 53, 52, 51, 184, 185},
		Over18:       true,
		Pager:        ptttype.PAGER_ON,
		Career:       ptttype.Career_t{165, 254, 180, 186, 179, 110, 197, 233},
		LastSeen:     1600681288,
		UaVersion:    128,
	}

	testNewRegister1Passwd = []byte("!@Ab86")

	testBoardSummary6 = &ptttype.BoardSummaryRaw{
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

	testBoardSummary7 = &ptttype.BoardSummaryRaw{
		Bid:     7,
		Brdname: &ptttype.BoardID_t{0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x00, 0x2e, 0x2e, 0x2e, 0x2e},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xb8, 0xea, 0xb7,
			0xbd, 0xa6, 0x5e, 0xa6, 0xac, 0xb5, 0xa9, 0x00, 0xb7, 0x73,
			0xa4, 0xe5, 0xb3, 0xb9, 0x00, 0x20, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x00, 0x6e,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardSummary11 = &ptttype.BoardSummaryRaw{
		Bid:     11,
		Brdname: &ptttype.BoardID_t{0x45, 0x64, 0x69, 0x74, 0x45, 0x78, 0x70, 0x00, 0x2e, 0x2e, 0x2e, 0x2e},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xbd, 0x64, 0xa5,
			0xbb, 0xba, 0xeb, 0xc6, 0x46, 0xa7, 0xeb, 0xbd, 0x5a, 0xb0,
			0xcf, 0x00, 0xd6, 0xa1, 0x49, 0x00, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x00, 0x6e,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardSummary8 = &ptttype.BoardSummaryRaw{
		Bid:     8,
		Brdname: &ptttype.BoardID_t{0x4e, 0x6f, 0x74, 0x65, 0x00, 0x65, 0x64, 0x00, 0x2e, 0x2e, 0x2e, 0x2e},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xb0, 0xca, 0xba,
			0x41, 0xac, 0xdd, 0xaa, 0x4f, 0xa4, 0xce, 0xba, 0x71, 0xa6,
			0xb1, 0xa7, 0xeb, 0xbd, 0x5a, 0x00, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x00, 0x6e,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardSummary1 = &ptttype.BoardSummaryRaw{
		Bid:     1,
		BrdAttr: ptttype.BRD_POSTMASK,
		Brdname: &ptttype.BoardID_t{'S', 'Y', 'S', 'O', 'P'},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xaf, 0xb8, 0xaa,
			0xf8, 0xa6, 0x6e, 0x21,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardSummary10 = &ptttype.BoardSummaryRaw{
		Bid:     10,
		Brdname: &ptttype.BoardID_t{'W', 'h', 'o', 'A', 'm', 'I', 0x00, 0x00, 0x2e, 0x2e, 0x2e, 0x2e},
		Title: &ptttype.BoardTitle_t{
			0xbc, 0x54, 0xad, 0xf9, 0x20, 0xa1, 0xb7, 0xa8, 0xfe, 0xa8,
			0xfe, 0xa1, 0x41, 0xb2, 0x71, 0xb2, 0x71, 0xa7, 0xda, 0xac,
			0x4f, 0xbd, 0xd6, 0xa1, 0x49, 0x00, 0xaf, 0xb8, 0xaa, 0xf8,
			0x20, 0x20, 0xa3, 0xad, 0xa1, 0x49, 0x0, 0x6e,
		},
		BM:       []*ptttype.UserID_t{},
		StatAttr: ptttype.NBRD_FAV,
	}

	testBoardHeaderRaw1 = &ptttype.BoardHeaderRaw{
		Brdname: ptttype.BoardID_t{'t', 'e', 's', 't'},
		BrdAttr: ptttype.BRD_HIDE,
	}
)
