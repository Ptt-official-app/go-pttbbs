package cmsys

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func GetNumRecords(filename string, size uintptr) int {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0
	}

	return int(stat.Size() / int64(size))
}

func GetRecords(boardID *ptttype.BoardID_t, filename string, startAid ptttype.Aid, n int, isDesc bool, maxAid ptttype.Aid) (summaries []*ptttype.ArticleSummaryRaw, err error) {
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

	//construct headers
	summaries = make([]*ptttype.ArticleSummaryRaw, 0, n)
	logrus.Infof("cmsys.GetRecords: startAid: %v maxAid: %v n: %v isDesc: %v", startAid, maxAid, n, isDesc)
	for idx, aid, aidInFile := 0, startAid, startAid.ToAidInStore(); idx < n; idx++ {
		logrus.Infof("cmsys.GetRecords: (%v/%v) aid: %v aidInFile: %v", idx, n, aid, aidInFile)
		if aid == 0 || aid > maxAid {
			break
		}

		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(aidInFile), 0)
		if err != nil {
			return summaries, nil
		}

		header := &ptttype.FileHeaderRaw{}
		err = binary.Read(file, binary.LittleEndian, header)
		if err != nil {
			return summaries, nil
		}

		each := ptttype.NewArticleSummaryRaw(aid, boardID, header)
		summaries = append(summaries, each)

		if isDesc {
			aid--
			aidInFile = aid.ToAidInStore()
		} else {
			aid++
			aidInFile = aid.ToAidInStore()
		}
	}

	return summaries, nil
}

//FindRecordStartAid
//
//startIdx should be 1-total.
func FindRecordStartAid(dirFilename string, total int, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (startAid ptttype.Aid, err error) {

	file, err := os.Open(dirFilename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return -1, err
	}
	defer file.Close()

	start := 0
	end := int(total) - 1

	//binary-search based on create-time.
	idxInStore := 0
	header := &ptttype.FileHeaderRaw{}
	for idxInStore = (start + end) / 2; ; idxInStore = (start + end) / 2 {
		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
		if err != nil {
			return -1, err
		}

		err = binary.Read(file, binary.LittleEndian, header)
		if err != nil {
			return -1, err
		}

		fileCreateTime, err := header.Filename.CreateTime()
		if err != nil {
			return -1, err
		}
		j := createTime - fileCreateTime

		logrus.Infof("cmsys.FindRecordStartAid: start: %v end: %v idxInStore: %v j: %v", start, end, idxInStore, j)
		if j == 0 {
			break
		}

		if end == start {
			break
		} else if idxInStore == start {
			idxInStore = end
			start = end
		} else if j > 0 {
			start = idxInStore
		} else {
			end = idxInStore
		}
	}

	fileCreateTime, _ := header.Filename.CreateTime()
	if createTime == fileCreateTime && filename != nil && bytes.Equal(filename[:], header.Filename[:]) {
		return ptttype.Aid(idxInStore + 1), nil
	}

	//find the start
	logrus.Infof("cmsys.FindRecordStartAid: to find start: createTime: %v fileCreateTime: %v idxInStore: %v total: %v isDesc: %v", createTime, fileCreateTime, idxInStore, total, isDesc)
	if isDesc {
		for ; idxInStore < total; idxInStore++ {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = binary.Read(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime < fileCreateTime {
				break
			}

		}
		if idxInStore == total {
			idxInStore = total - 1
		}
	} else {
		for ; idxInStore >= 0; idxInStore-- {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = binary.Read(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime > fileCreateTime {
				break
			}

		}
		if idxInStore == -1 {
			idxInStore = 0
		}

	}

	//linear search
	if isDesc {
		// it's supposed that createTime <= fileCreateTime for now.
		for ; idxInStore >= 0; idxInStore-- {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = binary.Read(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime == fileCreateTime && filename != nil && bytes.Equal(filename[:], header.Filename[:]) {
				return ptttype.Aid(idxInStore + 1), nil
			} else if createTime > fileCreateTime {
				break
			}
		}
	} else {
		// it's supposed that createTime >= fileCreateTime for now.
		for ; idxInStore < total; idxInStore++ {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = binary.Read(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime == fileCreateTime && filename != nil && bytes.Equal(filename[:], header.Filename[:]) {
				return ptttype.Aid(idxInStore + 1), nil
			} else if createTime < fileCreateTime {
				break
			}
		}
		if idxInStore == total {
			idxInStore = -1
		}
	}
	if idxInStore == -1 {
		return -1, nil
	}

	return ptttype.Aid(idxInStore + 1), nil
}
