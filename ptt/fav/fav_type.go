package fav

type FavType struct {
	TheType FavT
	Attr    Favh
	Fp      interface{}
}

func (ft *FavType) castFolder() *FavFolder {
	if ft.Fp != nil {
		favFolder, ok := ft.Fp.(*FavFolder)
		if !ok {
			return nil
		}
		return favFolder
	}
	favFolder := &FavFolder{}
	ft.Fp = favFolder
	return favFolder
}

func (ft *FavType) castBoard() *FavBoard {
	if ft.Fp != nil {
		favBoard, ok := ft.Fp.(*FavBoard)
		if !ok {
			return nil
		}
		return favBoard
	}

	favBoard := &FavBoard{}
	ft.Fp = favBoard
	return favBoard
}

func (ft *FavType) castLine() *FavLine {
	if ft.Fp != nil {
		favLine, ok := ft.Fp.(*FavLine)
		if !ok {
			return nil
		}
		return favLine
	}

	favLine := &FavLine{}
	ft.Fp = favLine
	return favLine
}

func (ft *FavType) isValid() bool {
	return ft.Attr&FAVH_FAV != 0
}

func (ft *FavType) copyTo(ft1 *FavType) {
	ft1.TheType = ft.TheType
	ft1.Attr = ft.Attr
	ft1.Fp = ft.Fp
}

func (ft *FavType) GetID() (theID int) {
	switch ft.TheType {
	case FAVT_BOARD:
		fp := ft.castBoard()
		return int(fp.Bid)
	case FAVT_LINE:
		fp := ft.castLine()
		return int(fp.Lid)
	case FAVT_FOLDER:
		fp := ft.castFolder()
		return int(fp.Fid)
	}

	return 0
}
