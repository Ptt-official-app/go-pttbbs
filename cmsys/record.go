package cmsys

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func GetNumRecords(filename string, size uintptr) int {
	stat, err := os.Stat(filename)
	if err != nil {
		return 0
	}

	return int(stat.Size() / int64(size))
}

func GetRecords(boardID *ptttype.BoardID_t, filename string, startIdx ptttype.SortIdx, n int, isDesc bool) (summaries []*ptttype.ArticleSummaryRaw, err error) {
	if !startIdx.IsValid() {
		return nil, ptttype.ErrInvalidIdx
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := info.Size()
	maxIdx := ptttype.SortIdx(fileSize / int64(ptttype.FILE_HEADER_RAW_SZ))

	// construct headers
	summaries = make([]*ptttype.ArticleSummaryRaw, 0, n)
	for i, idx, idxInFile := 0, startIdx, startIdx.ToSortIdxInStore(); i < n; i++ {
		if idx == 0 || idx > maxIdx {
			break
		}

		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInFile), 0)
		if err != nil {
			return summaries, nil
		}

		header := &ptttype.FileHeaderRaw{}
		err = types.BinaryRead(file, binary.LittleEndian, header)
		if err != nil {
			return summaries, nil
		}

		each := ptttype.NewArticleSummaryRaw(idx, boardID, header)
		summaries = append(summaries, each)

		if isDesc {
			idx--
			idxInFile = idx.ToSortIdxInStore()
		} else {
			idx++
			idxInFile = idx.ToSortIdxInStore()
		}
	}

	return summaries, nil
}

func GetRecord(dirFilename string, filename *ptttype.Filename_t, total int) (idx ptttype.SortIdx, fhdr *ptttype.FileHeaderRaw, err error) {
	createTime, err := filename.CreateTime()
	if err != nil {
		return 0, nil, err
	}
	idx, err = FindRecordStartIdx(dirFilename, total, createTime, filename, true)
	if err != nil {
		return 0, nil, err
	}

	file, err := os.Open(dirFilename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return 0, nil, err
	}
	defer file.Close()

	idxInFile := idx.ToSortIdxInStore()
	_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInFile), os.SEEK_SET)
	if err != nil {
		return 0, nil, err
	}

	fhdr = &ptttype.FileHeaderRaw{}
	err = types.BinaryRead(file, binary.LittleEndian, fhdr)
	if err != nil {
		return 0, nil, err
	}

	return idx, fhdr, nil
}

//FindRecordStartIdx
//
//startIdx should be 1-total.
func FindRecordStartIdx(dirFilename string, total int, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (startIdx ptttype.SortIdx, err error) {
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

	// binary-search based on create-time.
	idxInStore := 0
	header := &ptttype.FileHeaderRaw{}
	for idxInStore = (start + end) / 2; ; idxInStore = (start + end) / 2 {
		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
		if err != nil {
			return -1, err
		}

		err = types.BinaryRead(file, binary.LittleEndian, header)
		if err != nil {
			return -1, err
		}

		fileCreateTime, err := header.Filename.CreateTime()
		if err != nil {
			return -1, err
		}
		j := createTime - fileCreateTime

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
		return ptttype.SortIdx(idxInStore + 1), nil
	}

	// find the start
	if isDesc {
		for ; idxInStore < total; idxInStore++ {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = types.BinaryRead(file, binary.LittleEndian, header)
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

			err = types.BinaryRead(file, binary.LittleEndian, header)
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

	// linear search
	if isDesc {
		// it's supposed that createTime <= fileCreateTime for now.
		for ; idxInStore >= 0; idxInStore-- {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, err
			}

			err = types.BinaryRead(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime == fileCreateTime && filename != nil && bytes.Equal(filename[:], header.Filename[:]) {
				return ptttype.SortIdx(idxInStore + 1), nil
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

			err = types.BinaryRead(file, binary.LittleEndian, header)
			if err != nil {
				return -1, err
			}
			fileCreateTime, _ = header.Filename.CreateTime()

			if createTime == fileCreateTime && filename != nil && bytes.Equal(filename[:], header.Filename[:]) {
				return ptttype.SortIdx(idxInStore + 1), nil
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

	return ptttype.SortIdx(idxInStore + 1), nil
}

func SubstituteRecord(filename string, data interface{}, theSize uintptr, idxInStore int32) (err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(ptttype.DEFAULT_FILE_CREATE_PERM))
	if err != nil {
		return err
	}
	defer file.Close()

	offset := int64(idxInStore) * int64(theSize)
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	err = GoPttLock(file, filename, offset, theSize)
	if err != nil {
		return err
	}
	defer GoPttUnlock(file, filename, offset, theSize)

	err = types.BinaryWrite(file, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	return nil
}

func AppendRecord(filename string, data interface{}, theSize uintptr) (idx ptttype.SortIdx, err error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(ptttype.DEFAULT_FILE_CREATE_PERM))
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fd := file.Fd()
	err = GoFlock(fd, filename)
	if err != nil {
		return 0, err
	}
	defer GoFunlock(fd, filename)

	fsize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	idxInStore := ptttype.SortIdxInStore(fsize / int64(theSize))
	offset := int64(idxInStore) * int64(theSize)
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	err = types.BinaryWrite(file, binary.LittleEndian, data)
	if err != nil {
		return 0, err
	}

	return idxInStore.ToSortIdx(), nil
}
