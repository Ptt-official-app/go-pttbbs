package fav

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

/*
test0:
[35, 13,
 3, 0, 2, 1,
 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
 3, 1, 1,
 2, 1, 1, 183, 115, 170, 186, 165, 216, 191, 253, 0, 0,
 	   0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
 	   0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
 	   0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
       0, 0, 0, 0, 0, 0, 0, 0, 0,
 3, 1, 2,
 1, 1, 9, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
 1, 1, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,

 1, 0, 0, 0,
 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
]
*/

var (
	testTitle0 = [ptttype.BTLEN + 1]byte{
		183, 115, 170, 186, 165, 216, 191, 253, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0,
	}
	testSubFav0 = &FavRaw{
		FavNum:   1,
		NBoards:  1,
		NLines:   0,
		NFolders: 0,
		LineID:   0,
		FolderID: 0,
		Favh: []*FavType{
			{FAVT_BOARD, 1, &FavBoard{1, 0, 0}},
		},
	}

	testFav0 = &FavRaw{
		NBoards:  3,
		NLines:   2,
		NFolders: 1,
		LineID:   2,
		FolderID: 1,
		FavNum:   7,
		Favh: []*FavType{
			{FAVT_BOARD, 1, &FavBoard{1, 0, 0}},
			{FAVT_LINE, 1, &FavLine{1}},
			{FAVT_FOLDER, 1, &FavFolder{1, testTitle0, testSubFav0}},
			{FAVT_LINE, 1, &FavLine{2}},
			{FAVT_BOARD, 1, &FavBoard{9, 0, 0}},
			{FAVT_BOARD, 1, &FavBoard{8, 0, 0}},
		},
	}

	testFav1 = &FavRaw{
		FavNum:   6,
		NBoards:  2,
		NLines:   2,
		NFolders: 1,
		LineID:   2,
		FolderID: 1,
		Favh: []*FavType{
			{FAVT_LINE, 1, &FavLine{1}},
			{FAVT_FOLDER, 1, &FavFolder{1, testTitle0, testSubFav0}},
			{FAVT_LINE, 1, &FavLine{2}},
			{FAVT_BOARD, 1, &FavBoard{9, 0, 0}},
			{FAVT_BOARD, 1, &FavBoard{8, 0, 0}},
		},
	}
)
