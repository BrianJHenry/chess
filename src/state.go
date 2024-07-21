package chess

type Turn bool

const (
	WhiteTurn Turn = false
	BlackTurn Turn = true
)

type State struct {
	Board                   Board
	WhiteCanCastleKingSide  bool
	WhiteCanCastleQueenSide bool
	BlackCanCastleKingSide  bool
	BlackCanCastleQueenSide bool
	Turn                    Turn
	Previous                EncodedMove
}
