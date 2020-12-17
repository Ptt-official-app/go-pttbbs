package api

import (
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

var (
	testBoardSummary6 = &bbs.BoardSummary{
		BBoardID:   bbs.BBoardID("6_ALLPOST"),
		BrdAttr:    ptttype.BRD_POSTMASK,
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "ALLPOST",
		BoardClass: "嘰哩",
		RealTitle:  "跨板式LOCAL新文章",
		BoardType:  "◎",
		BM:         []bbs.UUserID{},
	}
	testBoardSummary7 = &bbs.BoardSummary{
		BBoardID:   bbs.BBoardID("7_deleted"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "deleted",
		BoardClass: "嘰哩",
		RealTitle:  "資源回收筒",
		BoardType:  "◎",
		BM:         []bbs.UUserID{},
	}
	testBoardSummary11 = &bbs.BoardSummary{
		BBoardID:   bbs.BBoardID("11_EditExp"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "EditExp",
		BoardClass: "嘰哩",
		RealTitle:  "範本精靈投稿區",
		BoardType:  "◎",
		BM:         []bbs.UUserID{},
	}
	testBoardSummary8 = &bbs.BoardSummary{
		BBoardID:   bbs.BBoardID("8_Note"),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "Note",
		BoardClass: "嘰哩",
		RealTitle:  "動態看板及歌曲投稿",
		BoardType:  "◎",
		BM:         []bbs.UUserID{},
	}

	testArticleSummary0 = &bbs.ArticleSummary{
		BBoardID:   bbs.BBoardID("10_WhoAmI"),
		ArticleID:  "1Vo_M_CD",
		IsDeleted:  false,
		Filename:   "M.1607202239.A.30D",
		CreateTime: 1607202239,
		MTime:      1607202238,
		Owner:      bbs.UUserID("SYSOP"),
		Date:       "12/06",
		Title:      "[問題] 我是誰？～",
		URL:        "http://localhost/bbs/WhoAmI/M.1607202239.A.30D.html",
	}

	testArticleSummary1 = &bbs.ArticleSummary{
		BBoardID:   bbs.BBoardID("10_WhoAmI"),
		ArticleID:  "1Vo_f30D",
		IsDeleted:  false,
		Filename:   "M.1607203395.A.00D",
		CreateTime: 1607203395,
		MTime:      1607203394,
		Owner:      bbs.UUserID("SYSOP"),
		Date:       "12/06",
		Title:      "[心得] 然後呢？～",
		Filemode:   ptttype.FILE_MARKED,
		URL:        "http://localhost/bbs/WhoAmI/M.1607203395.A.00D.html",
	}
)
