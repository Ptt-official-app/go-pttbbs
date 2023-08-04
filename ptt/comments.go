package ptt

import (
	"os"
	"time"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/Ptt-official-app/go-pttbbs/types/ansi"
	"github.com/sirupsen/logrus"
)

func FormatCommentString(user *ptttype.UserecRaw, board *ptttype.BoardHeaderRaw, commentType ptttype.CommentType, content []byte, ip *ptttype.IPv4_t, from []byte) (comment []byte, err error) {
	isLogIP := board.BrdAttr.HasPerm(ptttype.BRD_IPLOGRECMD)
	maxlen := 78 -
		3 - // lead
		6 - // date
		1 - // space
		6 // time

	if isLogIP {
		maxlen -= 15 // ip
	}

	// tail

	nowTSStr := types.NowTS().CdateMdHM()
	tail := make([]byte, 0, len(nowTSStr)+ptttype.IPV4LEN+len(from)+3)
	if isLogIP {
		tail = append(tail, types.CstrToBytes(ip[:])...)
	}
	tail = append(tail, ' ')
	tail = append(tail, []byte(nowTSStr)...)

	isAlignCmt := board.BrdAttr.HasPerm(ptttype.BRD_ALIGNEDCMT)
	var userBytes []byte
	if isAlignCmt {
		userBytes = user.UserID[:]
	} else {
		userBytes = types.CstrToBytes(user.UserID[:])
	}

	if ptttype.OLDRECOMMEND {
		maxlen -= 2
		maxlen -= len(userBytes)
		maxlen -= len(content)
		if maxlen < 0 {
			logrus.Errorf("FormatCommentString (old) maxlen < 0: userBytes: %v content: %v maxlen: %v", len(userBytes), len(content), maxlen)
			maxlen = 0
		}

		prefix := ansi.ANSIColor("1;31") + "\xa1\xf7 " + ansi.ANSIColor("33")
		infix := ansi.ANSIReset() + ansi.ANSIColor("33") + ":"
		whitespace := make([]byte, maxlen)
		for idx := range whitespace {
			whitespace[idx] = ' '
		}
		postfix := ansi.ANSIReset() + "\xb1\xc0"

		comment = make([]byte, 0, len(prefix)+len(userBytes)+len(infix)+len(content)+len(whitespace)+len(postfix)+len(tail)+1)
		comment = append(comment, []byte(prefix)...)
		comment = append(comment, userBytes...)
		comment = append(comment, []byte(infix)...)
		comment = append(comment, content...)
		comment = append(comment, whitespace...)
		comment = append(comment, []byte(postfix)...)
		comment = append(comment, tail...)
		comment = append(comment, '\n')

		return comment, nil
	}

	maxlen -= len(userBytes)
	maxlen -= len(content)
	if maxlen < 0 {
		logrus.Errorf("FormatCommentString (new) maxlen < 0: userBytes: %v content: %v maxlen: %v", len(userBytes), len(content), maxlen)
		maxlen = 0
	}

	commentTypeBytes := commentType.Bytes()
	prefix := ansi.ANSIColor("33")
	infix := ansi.ANSIReset() + ansi.ANSIColor("33") + ":"
	// XXX for ptt
	infix += " "
	whitespace := make([]byte, maxlen)
	for idx := range whitespace {
		whitespace[idx] = ' '
	}
	postfix := ansi.ANSIReset()

	comment = make([]byte, 0, len(commentTypeBytes)+1+len(prefix)+len(userBytes)+len(infix)+len(content)+len(whitespace)+len(postfix)+len(tail)+1)
	comment = append(comment, commentTypeBytes...)
	comment = append(comment, ' ')
	comment = append(comment, []byte(prefix)...)
	comment = append(comment, userBytes...)
	comment = append(comment, []byte(infix)...)
	comment = append(comment, content...)
	comment = append(comment, whitespace...)
	comment = append(comment, []byte(postfix)...)
	comment = append(comment, tail...)
	comment = append(comment, '\n')
	return comment, nil
}

const (
	DO_ADD_RECOMMEND_LOCK_WAIT = time.Duration(1 * time.Second)
)

func doAddRecommend(dirFilename string, idx ptttype.SortIdx, fhdr *ptttype.FileHeaderRaw, comment []byte, commentType ptttype.CommentType) (mtime types.Time4, err error) {
	filename := path.SetDIRPath(dirFilename, fhdr.Filename.String())

	if !ptttype.EDITPOST_SMARTMERGE {
		err = doAddRecommendNoSmartMerge(filename, comment)
		if err != nil {
			return 0, err
		}
	} else {
		for idxSmartMerge := 0; idxSmartMerge < 5; idxSmartMerge++ {
			err = doAddRecommendSmartMerge(filename, comment)
			if err == nil {
				break
			}
			time.Sleep(DO_ADD_RECOMMEND_LOCK_WAIT)
		}

		if err != nil {
			return 0, err
		}
	}

	update := int8(0)
	if commentType == ptttype.COMMENT_TYPE_RECOMMEND && fhdr.Recommend < ptttype.MAX_RECOMMENDS {
		update = 1
	} else if commentType == ptttype.COMMENT_TYPE_BOO && fhdr.Recommend > -ptttype.MAX_RECOMMENDS {
		update = -1
	}
	fhdr.Recommend += update
	fhdr.Modified = types.DashT(filename)
	if fhdr.Modified > 0 {
		err = ModifyDirLite(dirFilename, idx, &fhdr.Filename, fhdr.Modified, nil, nil, nil, update, nil, 0, 0)

		if err != nil {
			return 0, err
		}
	}

	return fhdr.Modified, nil
}

func doAddRecommendNoSmartMerge(filename string, comment []byte) (err error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(comment)

	return err
}

func doAddRecommendSmartMerge(filename string, comment []byte) (err error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = cmsys.GoFlockExNb(file.Fd(), filename)
	if err != nil {
		return err
	}
	defer func() { _ = cmsys.GoFunlock(file.Fd(), filename) }()

	_, err = file.Write(comment)

	return err
}
