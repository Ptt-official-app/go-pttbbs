package ptttype

type ChessCode byte

const (
	CHESSCODE_NONE    ChessCode = 0
	CHESSCODE_FIVE    ChessCode = 1
	CHESSCODE_CCHESS  ChessCode = 2
	CHESSCODE_GO      ChessCode = 3
	CHESSCODE_REVERSI ChessCode = 4
	CHESSCODE_MAX     ChessCode = 4
)

func (c ChessCode) IsValid() bool {
	return c >= CHESSCODE_NONE && c <= CHESSCODE_MAX
}
