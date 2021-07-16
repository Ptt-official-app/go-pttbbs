package ptttype

var origBBSHOME = ""

func SetIsTest() {
	origBBSHOME = SetBBSHOME("./testcase")

	initVars()
}

func UnsetIsTest() {
	SetBBSHOME(origBBSHOME)

	freeTestVars()
}

func freeTestVars() {
	ALLOW_EMAIL_LIST = nil
	REJECT_EMAIL_LIST = nil
	ALLOW_EMAIL_LIST_UPDATE_TS = 0
	REJECT_EMAIL_LIST_UPDATE_TS = 0
}
