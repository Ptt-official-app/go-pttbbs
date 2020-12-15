package cmsys

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func GetNumRecords(filename string, size uintptr) int {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0
	}

	return int(stat.Size() / int64(size))
}

func GetRecords(boardID *ptttype.BoardID_t, filename string, startAid ptttype.Aid, n int) ([]*ptttype.ArticleSummaryRaw, error) {
	if !startAid.IsValid() {
		return nil, ptttype.ErrInvalidAid
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}
	defer file.Close()

	aidInFile := startAid.ToAidInStore()
	_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(aidInFile), 0)
	if err != nil {
		return nil, err
	}

	//try read bytes
	bin := make([]byte, uintptr(n)*ptttype.FILE_HEADER_RAW_SZ)
	nRead, err := file.Read(bin)
	if nRead == 0 {
		//unable to read any.
		if err == io.EOF {
			err = nil
		}
		return nil, err
	}

	//construct headers
	nHeader := uintptr(uintptr(nRead) / ptttype.FILE_HEADER_RAW_SZ)
	headers := make([]ptttype.FileHeaderRaw, nHeader)
	buf := bytes.NewBuffer(bin)
	err = binary.Read(buf, binary.LittleEndian, headers)
	if err != nil {
		return nil, err
	}
	headers_p := make([]*ptttype.ArticleSummaryRaw, len(headers))
	for idx := 0; idx < len(headers); idx++ {
		headers_p[idx] = &ptttype.ArticleSummaryRaw{
			Aid:           startAid + ptttype.Aid(idx),
			BoardID:       boardID,
			FileHeaderRaw: &headers[idx],
		}
	}

	return headers_p, nil
}
