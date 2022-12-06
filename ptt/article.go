package ptt

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/sirupsen/logrus"
)

func CheckPostPerm2(uid ptttype.UID, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return CheckModifyPerm(uid, user, bid, board)
}

func CheckModifyPerm(uid ptttype.UID, user *ptttype.UserecRaw, bid ptttype.Bid, board *ptttype.BoardHeaderRaw) (err error) {
	return postpermMsg(uid, user, bid, board)
}

func getFileHeader(boardID *ptttype.BoardID_t, bid ptttype.Bid, filename *ptttype.Filename_t) (idx ptttype.SortIdx, fileHeader *ptttype.FileHeaderRaw, err error) {
	dirFilename, err := setBDir(boardID)
	if err != nil {
		return 0, nil, err
	}

	total, err := cache.GetBTotalWithRetry(bid)
	if err != nil {
		return 0, nil, err
	}
	if total == 0 {
		return 0, nil, ptttype.ErrInvalidFilename
	}

	idx, fileHeader, err = cmsys.GetRecord(dirFilename, filename, int(total))
	if err != nil {
		return 0, nil, err
	}

	return idx, fileHeader, nil
}

func ModifyDirLite(dirFilename string, idx ptttype.SortIdx, filename *ptttype.Filename_t, mtime types.Time4, title *ptttype.Title_t, owner *ptttype.Owner_t, theDate *ptttype.Date_t, recommend int8, multi []byte, enableModes ptttype.FileMode, disableModes ptttype.FileMode) (err error) {
	sz := types.DashS(dirFilename)
	if sz < int64(ptttype.FILE_HEADER_RAW_SZ)*int64(idx) {
		return ptttype.ErrInvalidIdx
	}

	file, err := os.OpenFile(dirFilename, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	idxInFile := idx.ToSortIdxInStore()
	_, err = file.Seek(int64(idxInFile)*int64(ptttype.FILE_HEADER_RAW_SZ), io.SeekStart)
	if err != nil {
		return err
	}

	fhdr := &ptttype.FileHeaderRaw{}
	err = types.BinaryRead(file, binary.LittleEndian, fhdr)
	if err != nil {
		return err
	}
	if types.Cstrcmp(fhdr.Filename[:], filename[:]) != 0 {
		return ptttype.ErrInvalidIdx
	}

	if mtime > 0 {
		fhdr.Modified = mtime
	}
	if enableModes != 0 {
		fhdr.Filemode |= enableModes
	}
	if disableModes != 0 {
		fhdr.Filemode &= ^disableModes
	}
	if title != nil && title[0] != 0 {
		fhdr.Title = *title
	}
	if owner != nil && owner[0] != 0 {
		fhdr.Owner = *owner
	}
	if theDate != nil && theDate[0] != 0 {
		fhdr.Date = *theDate
	}
	if multi != nil {
		copy(fhdr.Multi[:], multi)
	}

	if recommend != 0 {
		recommend += fhdr.Recommend
		if recommend > ptttype.MAX_RECOMMENDS {
			recommend = ptttype.MAX_RECOMMENDS
		} else if recommend < -ptttype.MAX_RECOMMENDS {
			recommend = -ptttype.MAX_RECOMMENDS
		}
		fhdr.Recommend = recommend
	}

	_, err = file.Seek(int64(idxInFile)*int64(ptttype.FILE_HEADER_RAW_SZ), io.SeekStart)
	if err != nil {
		return err
	}
	err = types.BinaryWrite(file, binary.LittleEndian, fhdr)

	return err
}

func readContent(filename string, retrieveTS types.Time4, isHash bool) (content []byte, mtime types.Time4, hash cmsys.Fnv64_t, err error) {
	// 1. check mtime
	stat, err := os.Stat(filename)
	if err != nil {
		return nil, 0, 0, err
	}
	mtime = types.TimeToTime4(stat.ModTime())
	if mtime <= retrieveTS {
		return nil, mtime, 0, nil
	}

	// 2. read content
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	content, err = io.ReadAll(file)
	if err != nil {
		return nil, 0, 0, err
	}

	if !isHash {
		return content, mtime, 0, nil
	}

	hash = cmsys.Fnv64Buf(content, len(content), cmsys.FNV1_64_INIT)
	// XXX do not do brc for now.
	// brcAddList(boardID, filename, updateTS)

	return content, mtime, hash, nil
}

// isFileOwner
//
// https://github.com/ptt/pttbbs/blob/master/mbbsd/bbs.c#L44
// XXX there are two known issues here:
// 1. Anonymous post don't work anymore.
// 2. People cross-posting (^X) very old post can't own it.
// Since ptt.cc does not have anonymous boards anymore, these issues are in
// low priority.  Sites using anonymous boards can fix on your own.
func isFileOwner(fhdr *ptttype.FileHeaderRaw, user *ptttype.UserecRaw) bool {
	if types.Cstrcmp(fhdr.Owner[:], user.UserID[:]) != 0 {
		return false
	}

	if types.Cstrlen(fhdr.Filename[:]) <= 3 {
		return false
	}

	createTS, _ := fhdr.Filename.CreateTime()

	return createTS >= user.FirstLogin
}

func hashPartialFile(filename string, sz int) (fnvseed cmsys.Fnv64_t, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fnvseed = cmsys.FNV1_64_INIT
	buf := make([]byte, 1024)
	reader := bufio.NewReader(file)
	for sz > 0 {
		nBuf, err := reader.Read(buf)
		if err == io.EOF {
			err = nil
		}
		if err != nil {
			return 0, err
		}
		if nBuf == 0 {
			break
		}
		if nBuf > sz {
			logrus.Warnf("nBuf > sz: nBuf: %v sz: %v", nBuf, sz)
			nBuf = sz
		}

		fnvseed = cmsys.Fnv64Buf(buf, nBuf, fnvseed)
		sz -= nBuf
	}
	if sz > 0 {
		return 0, ErrInvalidFileHash
	}

	return fnvseed, nil
}
