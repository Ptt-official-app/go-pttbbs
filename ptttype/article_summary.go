package ptttype

type ArticleSummaryRaw struct {
	Aid     Aid
	BoardID *BoardID_t

	*FileHeaderRaw

	Class []byte
}

func NewArticleSummaryRaw(idx SortIdx, boardID *BoardID_t, header *FileHeaderRaw) *ArticleSummaryRaw {

	class := header.Title.ToClass()

	return &ArticleSummaryRaw{
		Aid:     Aid(idx),
		BoardID: boardID,

		FileHeaderRaw: header,

		Class: class,
	}
}
