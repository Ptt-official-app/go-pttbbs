package ptt

import (
	"fmt"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

//WriteFile
//
//https://github.com/ptt/pttbbs/blob/master/mbbsd/edit.c#L3733
//https://github.com/ptt/pttbbs/blob/master/mbbsd/edit.c#L1924
func WriteFile(fpath string, flags ptttype.EditFlag, isSaveHeader bool, isUseAnony bool, title []byte, content [][]byte, user *ptttype.UserecRaw, uid ptttype.Uid, board *ptttype.BoardHeaderRaw, bid ptttype.Bid, ip *ptttype.IPv4_t, from []byte) (entropy int, err error) {

	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, ptttype.DEFAULT_FILE_CREATE_PERM)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	if isSaveHeader {
		_ = writeHeader(file, flags, title, user, board)
	}

	for idx, line := range content {
		if idx == len(content)-1 && len(line) == 0 { //ignore last empty line
			break
		}

		line = cmsys.Trim(line)
		line = StripANSIMoveCmd(line)

		fmt.Fprintf(file, "%s\n", line)
	}

	err = addSignature(file, isUseAnony, ip, from)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func writeHeader(file *os.File, flags ptttype.EditFlag, title []byte, user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) (err error) {

	err = writeHeaderAuthorBoard(file, flags, title, user, board)
	if err != nil {
		return err
	}

	nowTS := types.NowTS()
	fmt.Fprintf(file, "%s %s\n%s %s\n", ptttype.STR_TITLE_BIG5, title, ptttype.STR_TIME_BIG5, nowTS.Ctime())
	return nil

}

func writeHeaderAuthorBoard(file *os.File, flags ptttype.EditFlag, title []byte, user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) (err error) {
	if flags&(ptttype.EDITFLAG_KIND_MAILLIST|ptttype.EDITFLAG_KIND_SENDMAIL) != 0 {
		fmt.Fprintf(file, "%s %s (%s)\n", ptttype.STR_AUTHOR1_BIG5, types.CstrToBytes(user.UserID[:]), types.CstrToBytes(user.Nickname[:]))
		return nil
	}

	postLog := &PostLog{}
	author, nickname := writeHeaderAuthor(flags, user, board)
	copy(postLog.Author[:], author)
	copy(postLog.Board[:], types.CstrToBytes(board.Brdname[:]))
	copy(postLog.Title[:], title)
	nowTS := types.NowTS()
	postLog.TheDate = nowTS
	postLog.Number = 1

	_, err = cmsys.AppendRecord(path.SetBBSHomePath(ptttype.FN_POST), postLog, POSTLOG_SZ)
	if err != nil {
		return err
	}

	fmt.Fprintf(file, "%s %s (%s) %s %s\n", ptttype.STR_AUTHOR1_BIG5, author, nickname, ptttype.STR_POST1_BIG5, types.CstrToBytes(board.Brdname[:]))
	return nil
}

func writeHeaderAuthor(flags ptttype.EditFlag, user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw) (author []byte, nickname []byte) {
	author = types.CstrToBytes(user.UserID[:])
	nickname = types.CstrToBytes(user.Nickname[:])

	if !ptttype.HAVE_ANONYMOUS || board.BrdAttr&ptttype.BRD_ANONYMOUS == 0 {
		return author, nickname
	}

	author = ptttype.ANONYMOUS_ID_BYTES
	nickname = ptttype.ANONYMOUS_NICKNAME

	return author, nickname
}

//addSignature
//
//do not allow guest post.
//no need to do 簽名檔. already provided by user.
func addSignature(file *os.File, isUseAnony bool, ip *ptttype.IPv4_t, from []byte) (err error) {
	return addSimpleSignature(file, isUseAnony, ip, from)
}

func addSimpleSignature(file *os.File, isUseAnony bool, ip *ptttype.IPv4_t, from []byte) (err error) {
	var host []byte
	if isUseAnony {
		host = ptttype.ANONYMOUS_HOST
	} else {
		host = append(types.CstrToBytes(ip[:]), from...)
	}

	fmt.Fprintf(file, "\n--\n%s %s(%s), %s %s\n", ptttype.STR_BBS_BIG5, ptttype.BBSNAME_BIG5, ptttype.MYHOSTNAME, ptttype.STR_FROM_BIG5, host)

	return nil
}
