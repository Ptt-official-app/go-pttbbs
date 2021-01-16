package fav

import "github.com/Ptt-official-app/go-pttbbs/ptttype"

func setupTest() {
	ptttype.SetIsTest()
}

func teardownTest() {
	ptttype.UnsetIsTest()
}
