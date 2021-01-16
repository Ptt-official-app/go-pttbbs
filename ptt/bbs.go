package ptt

import (
	"io/ioutil"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//ReadPost
//
//pmore is replaced by frontend.
//We just need to return the whole content.
//We do not update brc here, because it requires lots of file-disk.
//require middlewares to handle user-read-article.
func ReadPost(user *ptttype.UserecRaw, uid ptttype.Uid, boardID *ptttype.BoardID_t, bid ptttype.Bid, filename *ptttype.Filename_t, retrieveTS types.Time4) (content []byte, mtime types.Time4, err error) {

	//1. check valid filename
	if filename[0] == 'L' || filename[0] == 0 {
		return nil, 0, ErrInvalidParams
	}

	cache.StatInc(ptttype.STAT_READPOST)

	//2. check perm.
	board, err := cache.GetBCache(bid)
	if err != nil {
		return nil, 0, err
	}

	statAttr := boardPermStat(user, uid, board, bid)
	if statAttr == ptttype.NBRD_INVALID {
		return nil, 0, ErrNotPermitted
	}

	//3. get filename
	theFilename, err := path.SetBFile(boardID, filename.String())
	if err != nil {
		return nil, 0, err
	}

	//4. check mtime
	stat, err := os.Stat(theFilename)
	if err != nil {
		return nil, 0, err
	}
	mtime = types.TimeToTime4(stat.ModTime())
	if mtime <= retrieveTS {
		return nil, mtime, nil
	}

	file, err := os.Open(theFilename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	content, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	//XXX do not do brc for now.
	//brcAddList(boardID, filename, updateTS)

	return content, mtime, nil
}
