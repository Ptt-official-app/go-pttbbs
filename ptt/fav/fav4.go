package fav

import (
	"encoding/binary"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func TryFav4Load(userID *ptttype.UserID_t, filename string) (favrec *FavRaw, err error) {
	oldFilename, err := path.SetHomeFile(userID, FAV4)
	if err != nil {
		return nil, err
	}
	if !types.IsRegularFile(oldFilename) {
		return nil, nil
	}

	file, err := os.Open(oldFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	favrec, err = fav4ReadFavrec(file)
	if err != nil {
		return nil, err
	}
	_, err = favrec.Save(userID)
	if err != nil {
		return nil, err
	}

	bakFilename, err := path.SetHomeFile(userID, FAV+".bak")
	if err != nil {
		return favrec, nil
	}
	// XXX copy new fav-filename to bak in pttbbs
	_ = types.CopyFile(filename, bakFilename)

	return favrec, nil
}

func fav4ReadFavrec(file *os.File) (*FavRaw, error) {
	favrec := NewFavRaw(nil)
	err := types.BinaryRead(file, binary.LittleEndian, &favrec.NBoards)
	if err != nil {
		return nil, ErrInvalidFav4Record
	}

	err = types.BinaryRead(file, binary.LittleEndian, &favrec.NLines)
	if err != nil {
		return nil, ErrInvalidFav4Record
	}

	err = types.BinaryRead(file, binary.LittleEndian, &favrec.NFolders)
	if err != nil {
		return nil, ErrInvalidFav4Record
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
			return nil, ErrInvalidFav4Record
		}
		err = types.BinaryRead(file, binary.LittleEndian, &ft.Attr)
		if err != nil {
			return nil, ErrInvalidFav4Record
		}

		switch ft.TheType {
		case FAVT_FOLDER:
			favFolder := ft.castFolder()
			if favFolder == nil {
				return nil, ErrInvalidFav4Record
			}
			err = types.BinRead(file, favFolder, ft.TheType.GetFav4TypeSize())
			if err != nil {
				return nil, ErrInvalidFav4Record
			}
		case FAVT_BOARD:
			favBoard := ft.castBoard()
			if favBoard == nil {
				return nil, ErrInvalidFav4Record
			}
			err = types.BinRead(file, favBoard, ft.TheType.GetFav4TypeSize())
			if err != nil {
				return nil, ErrInvalidFav4Record
			}
		case FAVT_LINE:
			favLine := ft.castLine()
			if favLine == nil {
				return nil, ErrInvalidFav4Record
			}
			err = types.BinRead(file, favLine, ft.TheType.GetFav4TypeSize())
			if err != nil {
				return nil, ErrInvalidFav4Record
			}
		}
	}

	favrec.FavNum = nFavh
	for _, ft := range favrec.Favh {
		switch ft.TheType {
		case FAVT_FOLDER:
			favFolder := ft.castFolder()
			if favFolder == nil {
				return nil, ErrInvalidFavFolder
			}
			p, err := fav4ReadFavrec(file)
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
