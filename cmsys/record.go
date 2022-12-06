package cmsys

import (
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

		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInFile), io.SeekStart)
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
	_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInFile), io.SeekStart)
	if err != nil {
		return 0, nil, err
	}

	fhdr = &ptttype.FileHeaderRaw{}
	err = types.BinaryRead(file, binary.LittleEndian, fhdr)
	if err != nil {
		return 0, nil, err
	}
	if !filename.Eq(&fhdr.Filename) {
		return 0, nil, ErrRecordNotFound
	}

	return idx, fhdr, nil
}

// FindRecordStartIdx
//
// startIdx should be 1-total.
// find record:
// if isDesc: search from the newest, until either find the filename,
//
//	or the record of the createTime
//
// else:      search from the oldest, until either find the filename,
//
//	or the record of the createTime
func FindRecordStartIdx(dirFilename string, total int, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (startIdx ptttype.SortIdx, err error) {
	file, err := os.Open(dirFilename)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return -1, err
	}
	defer file.Close()

	header := &ptttype.FileHeaderRaw{}

	startStart := ptttype.SortIdxInStore(0)
	endEnd := ptttype.SortIdxInStore(total) - 1
	startStart, _, err = findValidRecordIdxInStore(startStart, file, header, false, startStart, endEnd)
	if err != nil {
		return -1, err
	}
	endEnd, _, err = findValidRecordIdxInStore(endEnd, file, header, true, startStart, endEnd)

	idxInStore, fileFilename, err := findRecordStartIdxBinSearch(file, startStart, endEnd, createTime, filename, isDesc)
	if err != nil {
		return -1, err
	}
	fileCreateTime, err := fileFilename.CreateTime()
	if err != nil {
		return -1, err
	}
	if createTime == fileCreateTime && filename != nil && filename.Eq(fileFilename) {
		return idxInStore.ToSortIdx(), nil
	}

	if isDesc {
		idxInStore, err = findRecordStartIdxPostSearchDesc(idxInStore, file, startStart, endEnd, createTime, filename)
		if err != nil {
			return -1, err
		}
	} else { // is ascending
		idxInStore, err = findRecordStartIdxPostSearchAsc(idxInStore, file, startStart, endEnd, createTime, filename)
		if err != nil {
			return -1, err
		}
	}

	return idxInStore.ToSortIdx(), nil
}

func findValidRecordIdxInStore(idxInStore ptttype.SortIdxInStore, file *os.File, header *ptttype.FileHeaderRaw, isDesc bool, start ptttype.SortIdxInStore, end ptttype.SortIdxInStore) (newIdxInStore ptttype.SortIdxInStore, createTime types.Time4, err error) {
	if !isDesc {
		for ; idxInStore <= end; idxInStore++ {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, 0, err
			}

			err = types.BinaryRead(file, binary.LittleEndian, header)
			if err != nil {
				return -1, 0, err
			}

			createTime, err = header.Filename.CreateTime()
			if err == nil {
				return idxInStore, createTime, nil
			}
		}
	} else {
		for ; idxInStore >= start; idxInStore-- {
			_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
			if err != nil {
				return -1, 0, err
			}

			err = types.BinaryRead(file, binary.LittleEndian, header)
			if err != nil {
				return -1, 0, err
			}

			createTime, err = header.Filename.CreateTime()
			if err == nil {
				return idxInStore, createTime, nil
			}
		}
	}

	return -1, 0, ErrRecordNotFound
}

func findRecordStartIdxBinSearch(file *os.File, startStart ptttype.SortIdxInStore, endEnd ptttype.SortIdxInStore, createTime types.Time4, filename *ptttype.Filename_t, isDesc bool) (idxInStore ptttype.SortIdxInStore, fileFilename *ptttype.Filename_t, err error) {
	start := startStart
	end := endEnd

	// binary-search based on create-time.
	// start, end should always be with valid idx.
	header := &ptttype.FileHeaderRaw{}
	idxInStore = ptttype.SortIdxInStore(0)
	for idxInStore = (start + end) / 2; ; idxInStore = (start + end) / 2 {
		_, err = file.Seek(int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idxInStore), io.SeekStart)
		if err != nil {
			return -1, nil, err
		}
		err = types.BinaryRead(file, binary.LittleEndian, header)
		if err != nil {
			return -1, nil, err
		}
		fileCreateTime, err := header.Filename.CreateTime()
		if err != nil {
			if start == end {
				break
			}
			idxInStore, fileCreateTime, start, end, err = findRecordStartIdxBinSearchValidIdxInStore(idxInStore, file, header, start, end)
			if err != nil {
				return -1, nil, err
			}
		}

		// we should have valid idxInStore here.
		j := createTime - fileCreateTime
		if j == 0 {
			break
		}
		if end == start {
			break
		} else if idxInStore == start {
			idxInStore = end //nolint
			start = end
		} else if j > 0 {
			start = idxInStore
		} else {
			end = idxInStore
		}
	}

	// We should have valid header
	return idxInStore, &header.Filename, nil
}

func findRecordStartIdxBinSearchValidIdxInStore(idxInStore ptttype.SortIdxInStore, file *os.File, header *ptttype.FileHeaderRaw, start ptttype.SortIdxInStore, end ptttype.SortIdxInStore) (newIdxInStore ptttype.SortIdxInStore, newFileCreateTime types.Time4, newStart ptttype.SortIdxInStore, newEnd ptttype.SortIdxInStore, err error) {
	if idxInStore == start {
		newIdxInStore, newFileCreateTime, err = findValidRecordIdxInStore(idxInStore, file, header, false, start, end)
		if err != nil {
			return -1, 0, 0, 0, err
		}
		return newIdxInStore, newFileCreateTime, newIdxInStore, end, nil
	} else if idxInStore == end {
		newIdxInStore, newFileCreateTime, err = findValidRecordIdxInStore(idxInStore, file, header, true, start, end)
		if err != nil {
			return -1, 0, 0, 0, err
		}
		return newIdxInStore, newFileCreateTime, start, newIdxInStore, nil
	} else {
		newIdxInStore, newFileCreateTime, err = findValidRecordIdxInStore(idxInStore, file, header, false, start, end)
		if err != nil {
			return -1, 0, 0, 0, err
		}
		if newIdxInStore == end {
			newIdxInStore, newFileCreateTime, err = findValidRecordIdxInStore(idxInStore, file, header, true, start, end)
			if err != nil {
				return -1, 0, 0, 0, err
			}
		}

		return newIdxInStore, newFileCreateTime, start, end, nil
	}
}

func findRecordStartIdxPostSearchDesc(idxInStore ptttype.SortIdxInStore, file *os.File, startStart ptttype.SortIdxInStore, endEnd ptttype.SortIdxInStore, createTime types.Time4, filename *ptttype.Filename_t) (newIdxInStore ptttype.SortIdxInStore, err error) {
	// find the start
	header := &ptttype.FileHeaderRaw{}
	for ; idxInStore <= endEnd; idxInStore++ {
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
			continue
		}

		if fileCreateTime > createTime {
			break
		}

	}
	if idxInStore > endEnd {
		idxInStore = endEnd
	}

	newIdxInStore, err = findRecordStartIdxPostSearchDescLinearSearch(idxInStore, file, startStart, createTime, filename)
	if err == nil {
		return newIdxInStore, nil
	}

	return findRecordStartIdxPostSearchDescLinearSearch(idxInStore, file, startStart, createTime, nil)
}

func findRecordStartIdxPostSearchDescLinearSearch(idxInStore ptttype.SortIdxInStore, file *os.File, startStart ptttype.SortIdxInStore, createTime types.Time4, filename *ptttype.Filename_t) (newIdxInStore ptttype.SortIdxInStore, err error) {
	// linear search
	// it's supposed that  fileCreateTime > createTime for now.
	header := &ptttype.FileHeaderRaw{}
	for ; idxInStore >= startStart; idxInStore-- {
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
			continue
		}

		if createTime == fileCreateTime && (filename == nil || filename.Eq(&header.Filename)) {
			return idxInStore, nil
		} else if fileCreateTime < createTime {
			if filename != nil {
				return -1, ErrRecordNotFound
			}
			break
		}
	}
	if idxInStore < startStart {
		return -1, ErrRecordNotFound
	}

	return idxInStore, nil
}

func findRecordStartIdxPostSearchAsc(idxInStore ptttype.SortIdxInStore, file *os.File, startStart ptttype.SortIdxInStore, endEnd ptttype.SortIdxInStore, createTime types.Time4, filename *ptttype.Filename_t) (newIdxInStore ptttype.SortIdxInStore, err error) {
	// find the start
	header := &ptttype.FileHeaderRaw{}
	for ; idxInStore >= startStart; idxInStore-- {
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
			continue
		}

		if fileCreateTime > createTime {
			break
		}

	}
	if idxInStore < startStart {
		idxInStore = startStart
	}

	newIdxInStore, err = findRecordStartIdxPostSearchAscLinearSearch(idxInStore, file, endEnd, createTime, filename)
	if err == nil {
		return newIdxInStore, nil
	}

	return findRecordStartIdxPostSearchAscLinearSearch(idxInStore, file, endEnd, createTime, nil)
}

func findRecordStartIdxPostSearchAscLinearSearch(idxInStore ptttype.SortIdxInStore, file *os.File, endEnd ptttype.SortIdxInStore, createTime types.Time4, filename *ptttype.Filename_t) (newIdxInStore ptttype.SortIdxInStore, err error) {
	// linear search
	// it's supposed that fileCreateTime < createTime for now.
	header := &ptttype.FileHeaderRaw{}
	for ; idxInStore <= endEnd; idxInStore++ {
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
			continue
		}

		if createTime == fileCreateTime && (filename == nil || filename.Eq(&header.Filename)) {
			return idxInStore, nil
		} else if fileCreateTime > createTime {
			if filename != nil {
				return -1, ErrRecordNotFound
			}
			break
		}
	}
	if idxInStore > endEnd {
		return -1, ErrRecordNotFound
	}

	return idxInStore, nil
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
	defer func() { _ = GoPttUnlock(file, filename, offset, theSize) }()

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
	defer func() { _ = GoFunlock(fd, filename) }()

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

func DeleteRecord(filename string, index ptttype.SortIdxInStore, theSize uintptr) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(ptttype.DEFAULT_FILE_CREATE_PERM))
	if err != nil {
		return err
	}
	defer file.Close()

	fd := file.Fd()
	err = GoFlock(fd, filename)
	if err != nil {
		return err
	}
	defer func() { _ = GoFunlock(fd, filename) }()

	offset := int64(index) * int64(theSize)
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}
	err = types.BinaryWrite(file, binary.LittleEndian, []byte(ptttype.FN_SAFEDEL))
	if err != nil {
		return err
	}
	return nil
}
