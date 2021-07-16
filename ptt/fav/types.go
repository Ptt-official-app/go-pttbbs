package fav

import (
	"unsafe"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

type (
	Fid int8
	Lid int8
)

type Favh int8

const (
	FAVH_FAV     Favh = 1
	FAVH_TAG     Favh = 2
	FAVH_UNREAD  Favh = 4
	FAVH_ADM_TAG Favh = 8
)

type FavLine struct {
	Lid Lid
}

const SIZE_OF_FAV_LINE = unsafe.Sizeof(FavLine{})

type FavFolder struct {
	Fid        Fid
	Title      ptttype.BoardTitle_t
	ThisFolder *FavRaw
}

const SIZE_OF_FAV_FOLDER = unsafe.Sizeof(FavFolder{})

type FavBoard struct {
	Bid       ptttype.Bid
	LastVisit int32 /* UNUSED */
	Attr      Favh
}

const SIZE_OF_FAV_BOARD = unsafe.Sizeof(FavBoard{})

type Fav4Folder struct {
	Fid        Fid
	Title      ptttype.BoardTitle_t
	ThisFolder int32
}

const SIZE_OF_FAV4_FOLDER = unsafe.Sizeof(Fav4Folder{})

type Fav4Board struct {
	Bid       ptttype.Bid
	LastVisit int32
	Attr      Favh
}

const SIZE_OF_FAV4_BOARD = unsafe.Sizeof(Fav4Board{})
