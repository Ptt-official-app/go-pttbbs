package bbs

import (
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
}

func NewArticleSummaryFromRaw(bboardID BBoardID, articleSummaryRaw *ptttype.ArticleSummaryRaw) *ArticleSummary {

	filename := types.CstrToString(articleSummaryRaw.Filename[:])
	createTime, _ := articleSummaryRaw.Filename.CreateTime()

	ownerID := ToUUserID(articleSummaryRaw.Owner.ToUserID())
	articleSummary := &ArticleSummary{
		BBoardID:   bboardID,
		ArticleID:  ToArticleID(&articleSummaryRaw.Filename, ownerID),
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

	return articleSummary
}
