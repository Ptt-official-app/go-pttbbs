package ptttype

type FriendType int

const (
	FRIEND_OVERRIDE FriendType = iota
	FRIEND_REJECT
	FRIEND_ALOHA
	_ // #define FRIEND_POST     3	    // deprecated
	FRIEND_SPECIAL
	FRIEND_CANVOTE
	BOARD_WATER
	BOARD_VISIBLE
)

var FriendFile = []string{ // for friend_edit type.
	FN_OVERRIDES,
	FN_REJECT,
	FN_ALOHAED,
	"",
	"",
	FN_CANVOTE,
	FN_WATER,
	FN_VISIBLE,
}

func (f FriendType) Filename() string {
	return FriendFile[f]
}
