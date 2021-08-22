package bbs

import (
	"strconv"
	"strings"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type ArticleSummary struct {
	BBoardID   BBoardID         `json:"bid"`
	ArticleID  ArticleID        `json:"aid"`
	IsDeleted  bool             `json:"deleted"`
	Filename   string           `json:"filename"`
	CreateTime types.Time4      `json:"create_time"`
	MTime      types.Time4      `json:"modified"`
	Recommend  int8             `json:"recommend"`
	Owner      UUserID          `json:"owner"`
	FullTitle  []byte           `json:"full_title"`
	Money      int32            `json:"money"`
	Filemode   ptttype.FileMode `json:"mode"`
	Class      []byte           `json:"class"`
	Read       bool             `json:"read"`
	Idx        string           `json:"idx"`

	RealTitle   []byte              `json:"title"`
	SubjectType ptttype.SubjectType `json:"subject_type"`
}

func NewArticleSummaryFromRaw(bboardID BBoardID, articleSummaryRaw *ptttype.ArticleSummaryRaw) *ArticleSummary {
	filename := types.CstrToString(articleSummaryRaw.Filename[:])
	createTime, _ := articleSummaryRaw.Filename.CreateTime()

	theType, realTitleWithClass := cmbbs.SubjectEx(&articleSummaryRaw.Title)
	if articleSummaryRaw.Owner.IsCorpse() {
		theType = ptttype.SUBJECT_DELETED
	}
	theClass, realTitle := titleToClass(realTitleWithClass)

	ownerID := ToUUserID(articleSummaryRaw.Owner.ToUserID())
	articleSummary := &ArticleSummary{
		BBoardID:    bboardID,
		ArticleID:   ToArticleID(&articleSummaryRaw.Filename),
		IsDeleted:   articleSummaryRaw.IsDeleted(),
		Filename:    filename,
		CreateTime:  createTime,
		MTime:       articleSummaryRaw.Modified,
		Recommend:   articleSummaryRaw.Recommend,
		Owner:       ownerID,
		FullTitle:   types.CstrToBytes(articleSummaryRaw.Title[:]),
		Money:       articleSummaryRaw.Money(),
		Class:       theClass,
		Filemode:    articleSummaryRaw.Filemode,
		SubjectType: theType,
		RealTitle:   realTitle,
	}
	articleSummary.Idx = SerializeArticleIdxStr(articleSummary)

	return articleSummary
}

func titleToClass(title []byte) (theClass []byte, realTitle []byte) {
	// class
	if len(title) < 6 {
		return nil, title
	}

	if title[0] != '[' || title[5] != ']' {
		return nil, title
	}

	theClass = title[1:5]
	realTitle = title[6:]
	for {
		if len(realTitle) == 0 {
			return theClass, nil
		}
		if realTitle[0] != ' ' {
			break
		}
		realTitle = realTitle[1:]
	}

	return theClass, realTitle
}

func SerializeArticleIdxStr(summary *ArticleSummary) (idxStr string) {
	return strconv.Itoa(int(summary.CreateTime)) + "@" + string(summary.ArticleID)
}

func DeserializeArticleIdxStr(idxStr string) (createTime types.Time4, articleID ArticleID, err error) {
	theList := strings.Split(idxStr, "@")
	if len(theList) != 2 {
		return 0, "", ErrInvalidParams
	}

	createTime_i, err := strconv.Atoi(theList[0])
	if err != nil {
		return 0, "", err
	}
	createTime = types.Time4(createTime_i)

	articleID = ArticleID(theList[1])
	filename := articleID.ToRaw()
	createTimeFromArticleID, err := filename.CreateTime()
	if err != nil {
		return 0, "", err
	}
	if createTime != createTimeFromArticleID {
		return 0, "", ErrInvalidParams
	}

	return createTime, articleID, nil
}
