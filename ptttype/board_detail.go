package ptttype

import "github.com/Ptt-official-app/go-pttbbs/types"

type BoardDetailRaw struct {
	Bid          Bid
	LastPostTime types.Time4
	Total        int32
	*BoardHeaderRaw
}
