package bbs

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type ArticleSummary struct {
	ArticleID  ArticleID        `json:"aid"`
	IsDeleted  bool             `json:"deleted"`
	Filename   string           `json:"filename"`
	CreateTime types.Time4      `json:"create_time"`
	Mtime      types.Time4      `json:"mtime"`
	Recommend  int8             `json:"recommend"`
	Owner      string           `json:"owner"`
	Date       string           `json:"date"`
	Title      string           `json:"title"`
	Money      int32            `json:"money"`
	Filemode   ptttype.FileMode `json:"mode"`
}

func NewArticleSummaryFromRaw(articleSummaryRaw *ptttype.ArticleSummaryRaw) *ArticleSummary {

	articleSummary := &ArticleSummary{}

	articleSummary.ArticleID = ToArticleID(articleSummaryRaw.Aid, articleSummaryRaw.Filename)
	articleSummary.IsDeleted = articleSummaryRaw.IsDeleted()
	articleSummary.Filename = types.CstrToString(articleSummaryRaw.Filename[:])
	articleSummary.CreateTime, _ = articleSummaryRaw.Filename.CreateTime()
	articleSummary.Mtime = articleSummaryRaw.Modified
	articleSummary.Recommend = articleSummaryRaw.Recommend
	articleSummary.Owner = types.CstrToString(articleSummaryRaw.Owner.ToUserID()[:])
	articleSummary.Date = types.CstrToString(articleSummaryRaw.Date[:])
	articleSummary.Title = types.Big5ToUtf8(articleSummaryRaw.Title[:])
	articleSummary.Money = articleSummaryRaw.Money()
	articleSummary.Filemode = articleSummaryRaw.Filemode

	return articleSummary
}
