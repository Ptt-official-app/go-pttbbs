package ptttype

var (
	testUserecRaw1 = &UserecRaw{
		Version:    PASSWD_VERSION,
		UserID:     UserID_t{0x53, 0x59, 0x53, 0x4f, 0x50}, // SYSOP
		PasswdHash: Passwd_t{33, 2, 99, 49, 88, 98, 116, 52, 107, 105, 119, 69, 69, 0},

		UFlag:        33557088,
		UserLevel:    7,
		NumLoginDays: 1,
		NumPosts:     0,
		FirstLogin:   1606686147,
		LastLogin:    1606701137,
		LastHost:     IPv4_t{49, 55, 50, 46, 49, 56, 46, 48, 46, 49},
		Over18:       true,
		Pager:        PAGER_ON,
	}
)
