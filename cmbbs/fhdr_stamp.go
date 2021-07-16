package cmbbs

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func fhdrStamp(boardFilename string, header *ptttype.FileHeaderRaw, theType ptttype.StampType) (fullFilename string, err error) {
	nowTS := types.NowTS()

	filename := ""
	switch theType {
	case ptttype.STAMP_FILE:
		for {
			rnd := rand.Intn(0xfff + 1)
			nowTS++
			filename = fmt.Sprintf("M.%d.A.%3.3X", nowTS, rnd)
			fullFilename = boardFilename + string(os.PathSeparator) + filename
			isValid, err := fhdrStampIsValidFilename(fullFilename)
			if err != nil {
				return "", err
			}
			if isValid {
				break
			}
		}
	case ptttype.STAMP_DIR:
		for {
			nowTS++
			filename = fmt.Sprintf("D%X", nowTS&0o7777)
			fullFilename = boardFilename + string(os.PathSeparator) + filename
			isValid, err := fhdrStampIsValidDir(fullFilename)
			if err != nil {
				return "", err
			}
			if isValid {
				break
			}
		}
	case ptttype.STAMP_LINK:
		for {
			nowTS++
			filename = fmt.Sprintf("S%X", nowTS)
			fullFilename = boardFilename + string(os.PathSeparator) + filename
			isValid, err := fhdrStampIsValidLink(fullFilename)
			if err != nil {
				return "", err
			}
			if isValid {
				break
			}
		}
	default:
		return "", ptttype.ErrInvalidType
	}

	copy(header.Filename[:], []byte(filename))
	theDate := nowTS.Cdatemd()
	copy(header.Date[:], []byte(theDate))

	return fullFilename, nil
}

func fhdrStampIsValidFilename(fullFilename string) (isValid bool, err error) {
	file, err := types.OpenCreate(fullFilename, os.O_WRONLY|os.O_EXCL)
	if os.IsExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer file.Close()

	return true, nil
}

func fhdrStampIsValidDir(fullFilename string) (isValid bool, err error) {
	err = types.Mkdir(fullFilename)
	if os.IsExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func fhdrStampIsValidLink(fullFilename string) (isValid bool, err error) {
	err = types.Symlink("temp", fullFilename)
	if os.IsExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func Stampfile(boardFilename string, header *ptttype.FileHeaderRaw) (filename string, err error) {
	*header = ptttype.EMPTY_FILE_HEADER_RAW

	return fhdrStamp(boardFilename, header, ptttype.STAMP_FILE)
}

func StampfileU(boardFilename string, header *ptttype.FileHeaderRaw) (filename string, err error) {
	return fhdrStamp(boardFilename, header, ptttype.STAMP_FILE)
}
