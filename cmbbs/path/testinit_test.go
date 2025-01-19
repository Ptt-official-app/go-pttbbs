package path

func setupTest() {
	SetIsTest()
}

func teardownTest() {
	defer UnsetIsTest()
}
