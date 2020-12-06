package types

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

var (
	MILLI_TS_TO_NANO_TS int64 = 1000000
)

func GetRandom() string {
	theUUID, _ := uuid.NewRandom()
	theBytes, _ := theUUID.MarshalBinary()

	return base64.RawURLEncoding.EncodeToString(theBytes)[:22]
}

func GetCurrentMilliTS() int64 {
	return time.Now().UnixNano() / MILLI_TS_TO_NANO_TS
}

func IsRegularFile(filename string) bool {
	// dashf in pttbbs
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fileInfo.Mode().IsRegular()
}

func ReadFile(filename string) ([]byte, error) {
	if !IsRegularFile(filename) {
		return nil, errors.New("not regular")
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//BinRead
//
//Read the file with identified size, the first bytes are fit in v.
//Assuming that the size of v is <= theSize.
//XXX not sure what to do if v > theSize.
func BinRead(file *os.File, v interface{}, theSize uintptr) error {
	vSize := binary.Size(v)
	if vSize > int(theSize) {
		log.WithFields(log.Fields{"theSize": theSize, "vSize": vSize, "type": reflect.TypeOf(v)}).Warn("binRead: theSize < vSize")
		return ErrBytesTooLarge
	}

	err := binary.Read(file, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	binSize := binary.Size(v)
	nToRead := int64(theSize) - int64(binSize)
	if nToRead < 0 {
		log.WithFields(log.Fields{"theSize": theSize, "binSize": binSize}).Warn("binRead: theSize < binSize")
		return nil
	}
	log.WithFields(log.Fields{"theSize": theSize, "binSize": binSize, "nToRead": nToRead}).Debug("to seek")
	_, err = file.Seek(nToRead, 1)
	if err != nil {
		return err
	}
	return nil
}

//BinWrite
//
//Write the file with identified size, the first bytes are with v.
//Assuming that the size of v is <= theSize.
//XXX not sure what to do if v > theSize.
func BinWrite(file *os.File, v interface{}, theSize uintptr) error {
	vSize := binary.Size(v)
	if vSize > int(theSize) {
		return ErrBytesTooLarge
	}
	err := binary.Write(file, binary.LittleEndian, v)
	if err != nil {
		return err
	}
	binSize := binary.Size(v)
	nToWrite := int64(theSize) - int64(binSize)
	if nToWrite < 0 {
		log.WithFields(log.Fields{"theSize": theSize, "binSize": binSize}).Warn("binWrite: theSize < binSize")
		return nil
	}
	log.WithFields(log.Fields{"theSize": theSize, "binSize": binSize, "nToWrite": nToWrite}).Debug("to write dummy")
	dummy := make([]byte, nToWrite)
	_, err = file.Write(dummy)
	if err != nil {
		return err
	}
	return nil
}

//ReadLine
//
//To avoid confusion, use only []byte is nil or not
//to know whether it's the end of file.
func ReadLine(reader *bufio.Reader) ([]byte, error) {
	if reader == nil {
		return nil, ErrNilReader
	}
	line, err := reader.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}

	if len(line) == 0 {
		return nil, io.EOF
	}

	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	return line, nil
}
