package fav

type FavT int8

const (
	FAVT_BOARD  FavT = 1
	FAVT_FOLDER FavT = 2
	FAVT_LINE   FavT = 3
)

func (t FavT) String() string {
	switch t {
	case FAVT_BOARD:
		return "Board"
	case FAVT_FOLDER:
		return "Folder"
	case FAVT_LINE:
		return "Line"
	default:
		return "unknown"
	}
}

func (t FavT) IsValidFavType() bool {
	switch t {
	case FAVT_BOARD:
		return true
	case FAVT_FOLDER:
		return true
	case FAVT_LINE:
		return true
	}
	return false
}

func (t FavT) GetTypeSize() uintptr {
	switch t {
	case FAVT_BOARD:
		return SIZE_OF_FAV_BOARD
	case FAVT_FOLDER:
		return SIZE_OF_FAV_FOLDER
	case FAVT_LINE:
		return SIZE_OF_FAV_LINE
	}
	return 0
}

func (t FavT) GetFav4TypeSize() uintptr {
	switch t {
	case FAVT_BOARD:
		return SIZE_OF_FAV4_BOARD
	case FAVT_FOLDER:
		return SIZE_OF_FAV_FOLDER
	case FAVT_LINE:
		return SIZE_OF_FAV_LINE
	}
	return 0
}
