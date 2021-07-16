package ptt

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cache"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptt/fav"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

func reginitFav(uid ptttype.Uid, user *ptttype.UserecRaw) (err error) {
	// XXX TODO

	file, err := os.Open(ptttype.FN_DEFAULT_FAVS)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	favrec := fav.NewFavRaw(nil)

	for line, err := types.ReadLine(reader); err == nil; line, err = types.ReadLine(reader) {
		// already chomped in readline
		lineStripped := cmsys.StripBlank(line)
		if len(lineStripped) == 0 || lineStripped[0] == '#' {
			continue
		}
		boardID := &ptttype.BoardID_t{}
		copy(boardID[:], []byte(line))
		bid, err := cache.GetBid(boardID)
		if err != nil {
			continue
		}
		if bid == 0 {
			continue
		}

		favrec.AddBoard(bid)
	}

	_, err = favrec.Save(&user.UserID)

	return err
}

func GetFavorites(userID *ptttype.UserID_t, retrieveTS types.Time4) (content []byte, mtime types.Time4, err error) {
	// 1. get filename
	filename, err := path.SetHomeFile(userID, fav.FAV)
	if err != nil {
		return nil, 0, err
	}

	// 2. check mtime
	mtime, err = getFavoritesGetMTime(userID, filename)
	if err != nil {
		return nil, 0, err
	}
	if mtime == 0 {
		return nil, 0, nil
	}
	if mtime <= retrieveTS {
		return nil, mtime, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	content, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	return content, mtime, nil
}

func getFavoritesGetMTime(userID *ptttype.UserID_t, filename string) (mtime types.Time4, err error) {
	stat, err := os.Stat(filename)
	if err == nil {
		mtime = types.TimeToTime4(stat.ModTime())
		return mtime, nil
	}

	if !os.IsNotExist(err) {
		return 0, err
	}
	favrec, err := fav.TryFav4Load(userID, filename)
	if err != nil {
		return 0, err
	}
	if favrec == nil {
		return 0, nil
	}

	stat, err = os.Stat(filename)
	if err != nil {
		return 0, err
	}

	mtime = types.TimeToTime4(stat.ModTime())
	return mtime, nil
}
