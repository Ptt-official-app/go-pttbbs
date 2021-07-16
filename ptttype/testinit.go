package ptttype

var origBBSHOME = ""

func SetIsTest() {
	origBBSHOME = SetBBSHOME("./testcase")

	initVars()
}

func UnsetIsTest() {
	SetBBSHOME(origBBSHOME)
}
