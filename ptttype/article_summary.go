package ptttype

type ArticleSummaryRaw struct {
	Aid     Aid
	BoardID *BoardID_t

	*FileHeaderRaw
}

func NewArticleSummaryRaw(idx SortIdx, boardID *BoardID_t, header *FileHeaderRaw) *ArticleSummaryRaw {
	return &ArticleSummaryRaw{
		Aid:     Aid(idx),
		BoardID: boardID,

		FileHeaderRaw: header,
	}
}
