package types

func setupTest() {
	SetIsTest()
}

func teardownTest() {
	UnsetIsTest()
}
