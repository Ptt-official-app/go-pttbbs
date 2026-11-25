package bbs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

func TestLoadBoards_negativeMax(t *testing.T) {
	setupTest()
	defer teardownTest()

	for _, nBoards := range []int{-1, -5} {
		_, _, err := LoadFullClassBoards("SYSOP", 1, nBoards)
		assert.Equal(t, ErrInvalidParams, err, "LoadFullClassBoards nBoards: %v", nBoards)

		_, _, err = LoadAutoCompleteBoards("SYSOP", "", nBoards, "", true)
		assert.Equal(t, ErrInvalidParams, err, "LoadAutoCompleteBoards nBoards: %v", nBoards)

		_, _, err = LoadGeneralBoardDetails("SYSOP", "", nBoards, true, ptttype.BSORT_BY_NAME)
		assert.Equal(t, ErrInvalidParams, err, "LoadGeneralBoardDetails nBoards: %v", nBoards)

		_, _, err = LoadGeneralBoards("SYSOP", "", nBoards, nil, nil, true, ptttype.BSORT_BY_NAME)
		assert.Equal(t, ErrInvalidParams, err, "LoadGeneralBoards nBoards: %v", nBoards)
	}
}
