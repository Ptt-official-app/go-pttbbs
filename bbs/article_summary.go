package bbs

import (
	"strconv"
	"strings"

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
	Title      []byte           `json:"title"`
	Money      int32            `json:"money"`
	Filemode   ptttype.FileMode `json:"mode"`
	Class      []byte           `json:"class"`
	Read       bool             `json:"read"`
	Idx        string           `json:"idx"`
}

func NewArticleSummaryFromRaw(bboardID BBoardID, articleSummaryRaw *ptttype.ArticleSummaryRaw) *ArticleSummary {

	filename := types.CstrToString(articleSummaryRaw.Filename[:])
	createTime, _ := articleSummaryRaw.Filename.CreateTime()

	ownerID := ToUUserID(articleSummaryRaw.Owner.ToUserID())
	articleSummary := &ArticleSummary{
		BBoardID:   bboardID,
		ArticleID:  ToArticleID(&articleSummaryRaw.Filename),
		IsDeleted:  articleSummaryRaw.IsDeleted(),
		Filename:   filename,
		CreateTime: createTime,
		MTime:      articleSummaryRaw.Modified,
		Recommend:  articleSummaryRaw.Recommend,
		Owner:      ownerID,
		Title:      types.CstrToBytes(articleSummaryRaw.Title[:]),
		Money:      articleSummaryRaw.Money(),
		Class:      articleSummaryRaw.Class,
		Filemode:   articleSummaryRaw.Filemode,
	}
	articleSummary.Idx = SerializeArticleIdxStr(articleSummary)

	return articleSummary
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
