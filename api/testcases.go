package api

import (
	"strconv"

	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

var (
	testBoardSummary6 = &bbs.BoardSummary{
		Bid:        strconv.Itoa(6),
		BrdAttr:    ptttype.BRD_POSTMASK,
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "ALLPOST",
		BoardClass: "嘰哩",
		RealTitle:  "跨板式LOCAL新文章",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary7 = &bbs.BoardSummary{
		Bid:        strconv.Itoa(7),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "deleted",
		BoardClass: "嘰哩",
		RealTitle:  "資源回收筒",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary11 = &bbs.BoardSummary{
		Bid:        strconv.Itoa(11),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "EditExp",
		BoardClass: "嘰哩",
		RealTitle:  "範本精靈投稿區",
		BoardType:  "◎",
		BM:         []string{},
	}
	testBoardSummary8 = &bbs.BoardSummary{
		Bid:        strconv.Itoa(8),
		StatAttr:   ptttype.NBRD_FAV,
		Brdname:    "Note",
		BoardClass: "嘰哩",
		RealTitle:  "動態看板及歌曲投稿",
		BoardType:  "◎",
		BM:         []string{},
	}
)
