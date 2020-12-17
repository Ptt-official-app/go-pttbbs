package bbs

import (
	"io"
	"log"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func OpenUserecFile(filename string) ([]*Userec, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ret := []*Userec{}

	for uidInFile := ptttype.UidInStore(0); ; uidInFile++ {
		uid := uidInFile.ToUid()
		user, eachErr := NewUserecWithFile(uid, file)
		if eachErr != nil {
			// io.EOF is reading correctly to the end the file.
			if eachErr == io.EOF {
				break
			}

			err = eachErr
			break
		}
		ret = append(ret, user)
	}

	return ret, err

}

func NewUserecWithFile(uid ptttype.Uid, file *os.File) (*Userec, error) {
	userecRaw, err := ptttype.NewUserecRawWithFile(file)
	if err != nil {
		return nil, err
	}

	user := NewUserecFromRaw(uid, userecRaw)

	return user, nil
}
