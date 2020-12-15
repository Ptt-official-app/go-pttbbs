package bbs

import (
	"fmt"

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
	Owner      string           `json:"owner"`
	Date       string           `json:"date"`
	Title      string           `json:"title"`
	Money      int32            `json:"money"`
	Filemode   ptttype.FileMode `json:"mode"`
	URL        string           `json:"url"`
}

func NewArticleSummaryFromRaw(bboardID BBoardID, articleSummaryRaw *ptttype.ArticleSummaryRaw) *ArticleSummary {

	filename := types.CstrToString(articleSummaryRaw.Filename[:])
	boardID := types.CstrToString(articleSummaryRaw.BoardID[:])
	createTime, _ := articleSummaryRaw.Filename.CreateTime()

	articleSummary := &ArticleSummary{
		BBoardID:   bboardID,
		ArticleID:  ToArticleID(articleSummaryRaw.Aid, articleSummaryRaw.Filename),
		IsDeleted:  articleSummaryRaw.IsDeleted(),
		Filename:   filename,
		CreateTime: createTime,
		MTime:      articleSummaryRaw.Modified,
		Recommend:  articleSummaryRaw.Recommend,
		Owner:      types.CstrToString(articleSummaryRaw.Owner.ToUserID()[:]),
		Date:       types.CstrToString(articleSummaryRaw.Date[:]),
		Title:      types.Big5ToUtf8(articleSummaryRaw.Title[:]),
		Money:      articleSummaryRaw.Money(),
		Filemode:   articleSummaryRaw.Filemode,
		URL:        fmt.Sprintf("%v/%v/%v.html", ptttype.URL_PREFIX, boardID, filename),
	}

	return articleSummary
}
