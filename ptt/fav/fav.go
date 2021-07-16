package fav

import (
	"encoding/binary"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//FavRaw
//
//It's with it's own serialize method and does not directly copy by struct.
//We can add MTime in FavRaw.
//The content of FavFolder
type FavRaw struct {
	MTime  types.Time4
	FavNum int16

	Root *FavRaw

	NBoards  int16 /* number of the boards */
	NLines   int8  /* number of the lines */
	NFolders int8  /* number of the folders */
	LineID   Lid   /* current max line id */
	FolderID Fid   /* current max folder id */

	Favh []*FavType
}

func NewFavRaw(root *FavRaw) (favRaw *FavRaw) {
	favRaw = &FavRaw{
		Root: root,
	}

	if root == nil {
		favRaw.Root = favRaw
	}

	return favRaw
}

//Load
//
//Load fav from file.
func Load(userID *ptttype.UserID_t) (favrec *FavRaw, err error) {
	filename, err := path.SetHomeFile(userID, FAV)
	if err != nil {
		return nil, err
	}

	if !types.IsRegularFile(filename) {
		favrec, err = TryFav4Load(userID, filename)
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			}
			return nil, err
		}

		return favrec, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	mtime := types.TimeToTime4(stat.ModTime())

	// version
	version := int16(0)
	err = types.BinaryRead(file, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}

	// favrec
	favrec, err = ReadFavrec(file)
	if err != nil {
		return nil, err
	}

	favrec.MTime = mtime

	return favrec, nil
}

//Save
//
//save fav to file.
//XXX use rename to reduce the probability of race-condition.
func (f *FavRaw) Save(userID *ptttype.UserID_t) (*FavRaw, error) {
	f.cleanup()

	filename, err := path.SetHomeFile(userID, FAV)
	if err != nil {
		return nil, err
	}

	// It's possible that the file does not exists.
	isToSave, newFav, err := f.checkIsToSave(userID, filename)
	if err != nil {
		return nil, err
	}
	if !isToSave {
		return newFav, nil
	}

	postfix := types.GetRandom()
	tmpFilename, err := path.SetHomeFile(userID, FAV+".tmp."+postfix)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(tmpFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	version := FAV_VERSION
	err = types.BinaryWrite(file, binary.LittleEndian, &version)
	if err != nil {
		return nil, err
	}

	err = f.WriteFavrec(file)
	if err != nil {
		return nil, err
	}

	// to rename
	err = os.Rename(tmpFilename, filename)
	if err != nil {
		return nil, err
	}

	// mtime is embedded in Load.
	newFav, err = Load(userID)
	if err != nil {
		return nil, err
	}

	return newFav, nil
}

func (f *FavRaw) checkIsToSave(userID *ptttype.UserID_t, filename string) (isToSave bool, newFav *FavRaw, err error) {
	stat, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil, nil
		}
		return false, nil, err
	}

	mtime := types.TimeToTime4(stat.ModTime())
	if mtime < f.MTime {
		return true, nil, nil
	} else if mtime == f.MTime {
		return false, f, nil
	} else {
		newFav, err = Load(userID)
		if err != nil {
			return false, nil, err
		}
		return false, newFav, nil
	}
}

func (f *FavRaw) getDataNumber() int16 {
	return f.NBoards + int16(f.NLines) + int16(f.NFolders)
}

func (f *FavRaw) cleanup() {
	if !f.isNeedRebuildFav() {
		return
	}

	f.rebuildFav()
}

func (f *FavRaw) isNeedRebuildFav() bool {
	if f.Favh == nil {
		return false
	}
	for _, ft := range f.Favh {
		if ft == nil {
			continue
		}

		if !ft.isValid() {
			return true
		}

		switch ft.TheType {
		case FAVT_BOARD:
		case FAVT_LINE:
		case FAVT_FOLDER:
			childFav := ft.castFolder().ThisFolder
			if childFav.isNeedRebuildFav() {
				return true
			}
		default:
			return true
		}
	}
	return false
}

/**
 * 清除 fp(dir) 中無效的 entry/dir。「無效」指的是沒有 FAVH_FAV flag，所以
 * 不包含不存在的看板。
 */
func (f *FavRaw) rebuildFav() {
	f.LineID = 0
	f.FolderID = 0
	f.NLines = 0
	f.NFolders = 0
	f.NBoards = 0

	if f.Favh == nil {
		return
	}

	j := 0
	for i, ft := range f.Favh {
		if !ft.isValid() {
			continue
		}

		switch ft.TheType {
		case FAVT_BOARD:
		case FAVT_LINE:
		case FAVT_FOLDER:
			childFav := ft.castFolder().ThisFolder
			childFav.rebuildFav()
		default:
			continue
		}

		f.increase(ft)
		if i != j {
			ft.copyTo(f.Favh[j])
		}
		j++
	}

	nFavh := f.getDataNumber()
	// to be consistent with the data-tail.
	f.Favh = f.Favh[:nFavh]
}

func (f *FavRaw) increase(ft *FavType) {
	switch ft.TheType {
	case FAVT_BOARD:
		f.NBoards++
	case FAVT_LINE:
		f.NLines++
		f.LineID++
		ftLine := ft.castLine()
		ftLine.Lid = f.LineID
	case FAVT_FOLDER:
		f.NFolders++
		f.FolderID++
		ftFolder := ft.castFolder()
		ftFolder.Fid = f.FolderID
	}
}

func (f *FavRaw) WriteFavrec(file *os.File) error {
	if f == nil {
		return nil
	}

	err := types.BinaryWrite(file, binary.LittleEndian, &f.NBoards)
	if err != nil {
		return err
	}
	err = types.BinaryWrite(file, binary.LittleEndian, &f.NLines)
	if err != nil {
		return err
	}
	err = types.BinaryWrite(file, binary.LittleEndian, &f.NFolders)
	if err != nil {
		return err
	}
	total := f.getDataNumber()

	for i := int16(0); i < total; i++ {
		ft := f.Favh[i]
		err := types.BinaryWrite(file, binary.LittleEndian, &ft.TheType)
		if err != nil {
			return err
		}
		err = types.BinaryWrite(file, binary.LittleEndian, &ft.Attr)
		if err != nil {
			return err
		}

		switch ft.TheType {
		case FAVT_FOLDER:
			favFolder := ft.castFolder()
			if favFolder == nil {
				return ErrInvalidFavFolder
			}
			err = types.BinaryWrite(file, binary.LittleEndian, &favFolder.Fid)
			if err != nil {
				return err
			}
			err = types.BinaryWrite(file, binary.LittleEndian, &favFolder.Title)
			if err != nil {
				return err
			}
		case FAVT_BOARD:
			favBoard := ft.castBoard()
			if favBoard == nil {
				return ErrInvalidFavBoard
			}
			err = types.BinWrite(file, favBoard, ft.TheType.GetTypeSize())
			if err != nil {
				return err
			}
		case FAVT_LINE:
			favLine := ft.castLine()
			if favLine == nil {
				return ErrInvalidFavLine
			}
			err = types.BinWrite(file, favLine, ft.TheType.GetTypeSize())
			if err != nil {
				return err
			}
		}
	}

	for i := int16(0); i < total; i++ {
		ft := f.Favh[i]
		if ft == nil {
			continue
		}
		if ft.TheType == FAVT_FOLDER {
			favFolder := ft.castFolder()
			if favFolder == nil {
				return ErrInvalidFavFolder
			}
			err := favFolder.ThisFolder.WriteFavrec(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadFavrec(file *os.File) (favrec *FavRaw, err error) {
	// https://github.com/ptt/pttbbs/blob/master/mbbsd/fav.c
	favrec = NewFavRaw(nil)
	err = types.BinaryRead(file, binary.LittleEndian, &favrec.NBoards)
	if err != nil {
		return nil, ErrInvalidFavRecord
	}

	err = types.BinaryRead(file, binary.LittleEndian, &favrec.NLines)
	if err != nil {
		return nil, ErrInvalidFavRecord
	}

	err = types.BinaryRead(file, binary.LittleEndian, &favrec.NFolders)
	if err != nil {
		return nil, ErrInvalidFavRecord
	}

	nFavh := favrec.getDataNumber()
	favrec.LineID = 0
	favrec.FolderID = 0
	favrec.Favh = make([]*FavType, nFavh)

	for i := int16(0); i < nFavh; i++ {
		ft := &FavType{}
		favrec.Favh[i] = ft

		err = types.BinaryRead(file, binary.LittleEndian, &ft.TheType)
		if err != nil {
			return nil, ErrInvalidFavType
		}
		if !ft.TheType.IsValidFavType() {
			return nil, ErrInvalidFavType
		}

		err = types.BinaryRead(file, binary.LittleEndian, &ft.Attr)
		if err != nil {
			return nil, ErrInvalidFavType
		}

		switch ft.TheType {
		case FAVT_FOLDER:
			favFolder := ft.castFolder()
			if favFolder == nil {
				return nil, ErrInvalidFavFolder
			}
			err = types.BinaryRead(file, binary.LittleEndian, &favFolder.Fid)
			if err != nil {
				return nil, ErrInvalidFavFolder
			}
			err = types.BinaryRead(file, binary.LittleEndian, &favFolder.Title)
			if err != nil {
				return nil, ErrInvalidFavFolder
			}
		case FAVT_BOARD:
			favBoard := ft.castBoard()
			if favBoard == nil {
				return nil, ErrInvalidFavBoard
			}
			err = types.BinRead(file, favBoard, ft.TheType.GetTypeSize())
			if err != nil {
				return nil, ErrInvalidFavBoard
			}
		case FAVT_LINE:
			favLine := ft.castLine()
			if favLine == nil {
				return nil, ErrInvalidFavLine
			}
			err = types.BinRead(file, favLine, ft.TheType.GetTypeSize())
			if err != nil {
				return nil, ErrInvalidFavLine
			}
		}
	}

	favrec.FavNum = nFavh
	for i := int16(0); i < nFavh; i++ {
		ft := favrec.Favh[i]
		if ft == nil {
			continue
		}

		switch ft.TheType {
		case FAVT_FOLDER:
			favFolder := ft.castFolder()
			if favFolder == nil {
				return nil, ErrInvalidFavFolder
			}
			p, err := ReadFavrec(file)
			if err != nil {
				return nil, ErrInvalidFavFolder
			}
			favFolder.ThisFolder = p
			favrec.FolderID++
			favFolder.Fid = favrec.FolderID

			p.Root = favrec.Root
			favrec.FavNum += p.FavNum
		case FAVT_LINE:
			favLine := ft.castLine()
			if favLine == nil {
				return nil, ErrInvalidFavLine
			}
			favrec.LineID++
			favLine.Lid = favrec.LineID
		}
	}

	return favrec, nil
}

func (f *FavRaw) AddBoard(bid ptttype.Bid) (favType *FavType, err error) {
	// check whether already added
	favBoard, err := f.GetBoard(bid)
	if err != nil {
		return nil, err
	}

	if favBoard != nil {
		return favBoard, nil
	}

	if f.isMaxSize() {
		return nil, ErrTooManyFavs
	}

	fp := &FavBoard{
		Bid: bid,
	}

	return f.PreAppend(FAVT_BOARD, fp)
}

func (f *FavRaw) GetBoard(bid ptttype.Bid) (favType *FavType, err error) {
	if !bid.IsValid() {
		return nil, ptttype.ErrInvalidBid
	}

	return f.GetFavItem(int(bid), FAVT_BOARD), nil
}

func (f *FavRaw) AddLine() (favType *FavType, err error) {
	if f.isMaxSize() {
		return nil, ErrTooManyFavs
	}

	if f.NLines >= MAX_LINE {
		return nil, ErrTooManyLines
	}

	fp := &FavLine{}
	return f.PreAppend(FAVT_LINE, fp)
}

func (f *FavRaw) AddFolder() (favType *FavType, err error) {
	if f.isMaxSize() {
		return nil, ErrTooManyFavs
	}

	if f.NLines >= MAX_FOLDER {
		return nil, ErrTooManyFolders
	}

	fp := &FavFolder{}
	fp.ThisFolder = NewFavRaw(f.Root)
	return f.PreAppend(FAVT_FOLDER, fp)
}

func (f *FavRaw) GetFavItem(theID int, theType FavT) (favType *FavType) {
	for _, each := range f.Favh {
		if each.TheType != theType {
			continue
		}

		eachID := each.GetID()
		if eachID == theID {
			return each
		}
	}

	return nil
}

//PreAppend
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/fav.c#L804
//Although it is named PreAppend, actually it appends to DataTail
func (f *FavRaw) PreAppend(theType FavT, fp interface{}) (favType *FavType, err error) {
	favt := &FavType{
		TheType: theType,
		Attr:    FAVH_FAV,
		Fp:      fp,
	}

	f.Favh = append(f.Favh, favt)
	f.Increase(theType, fp)

	return favt, nil
}

func (f *FavRaw) Increase(theType FavT, fp interface{}) {
	switch theType {
	case FAVT_BOARD:
		f.NBoards++
	case FAVT_LINE:
		f.NLines++
		f.LineID++
		ft, _ := fp.(*FavLine)
		ft.Lid = f.LineID
	case FAVT_FOLDER:
		f.NFolders++
		f.FolderID++
		ft, _ := fp.(*FavFolder)
		ft.Fid = f.FolderID
	}
	f.Root.FavNum++
}

func (f *FavRaw) isMaxSize() bool {
	return f.Root.FavNum >= MAX_FAV
}

func (f *FavRaw) CleanRoot() {
	f.Root = nil
	for _, each := range f.Favh {
		if each.TheType != FAVT_FOLDER {
			continue
		}

		ft, _ := each.Fp.(*FavFolder)
		ft.ThisFolder.CleanRoot()
	}
}
