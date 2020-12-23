package ptttype

type ArticleSummaryRaw struct {
	Aid     Aid
	BoardID *BoardID_t

	*FileHeaderRaw

	Class []byte
}

func NewArticleSummaryRaw(aid Aid, boardID *BoardID_t, header *FileHeaderRaw) *ArticleSummaryRaw {

	class := header.Title.ToClass()

	return &ArticleSummaryRaw{
		Aid:     aid,
		BoardID: boardID,

		FileHeaderRaw: header,

		Class: class,
	}
}
