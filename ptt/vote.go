package ptt

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func bSuckinfileInvis(file *os.File, fname string, board *ptttype.BoardHeaderRaw) (err error) {
	inFile, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	firstLine, err := types.ReadLine(reader)
	if err != nil {
		return err
	}

	postIdx := bytes.Index(firstLine, ptttype.STR_POST1_BIG5)
	if postIdx == -1 {
		postIdx = bytes.Index(firstLine, ptttype.STR_POST2_BIG5)
	}
	if postIdx != -1 {
		firstLinePrefix := firstLine[:postIdx]
		firstLine = append(firstLinePrefix, INVIS_BYTES_POST...)
	}
	_, err = file.Write(firstLine)
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	return err
}

func bSuckinfile(file *os.File, fname string) (err error) {
	inFile, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer inFile.Close()

	content, err := ioutil.ReadAll(inFile)
	if err != nil {
		return err
	}

	_, err = file.Write(content)

	return err
}
