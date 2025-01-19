package cache

import (
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"
)

type mapCount struct {
	count        int32
	updateTimeUS types.US
}

type mapTime struct {
	time         types.Time4
	updateTimeUS types.US
}

type Map struct {
	ReloadStartTimeUS types.US
	ReloadEndTimeUS   types.US

	BoardNBottom      map[ptttype.BoardID_t]mapCount
	BoardTotal        map[ptttype.BoardID_t]mapCount
	BoardLastPostTime map[ptttype.BoardID_t]mapTime
}

func NewMap() (retMap *Map, err error) {
	retMap = &Map{
		BoardNBottom:      make(map[ptttype.BoardID_t]mapCount),
		BoardTotal:        make(map[ptttype.BoardID_t]mapCount),
		BoardLastPostTime: make(map[ptttype.BoardID_t]mapTime),
	}

	return retMap, nil
}

func (m *Map) GetBNumber() (bnumber int32) {
	return int32(len(m.BoardTotal))
}

func (m *Map) Reset() {
	for k := range m.BoardNBottom {
		delete(m.BoardNBottom, k)
	}

	for k := range m.BoardTotal {
		delete(m.BoardTotal, k)
	}

	for k := range m.BoardLastPostTime {
		delete(m.BoardLastPostTime, k)
	}

	m.ReloadEndTimeUS = 0
	m.ReloadEndTimeUS = 0
}
