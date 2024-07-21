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

func (state State) ExecuteMove(move Move) State {
	whiteCanCastleKingSide := state.WhiteCanCastleKingSide
	whiteCanCastleQueenSide := state.WhiteCanCastleQueenSide
	blackCanCastleKingSide := state.BlackCanCastleKingSide
	blackCanCastleQueenSide := state.BlackCanCastleQueenSide

	if state.Turn == BlackTurn {
		if move.Start.Equals(Position{X: 4, Y: 7}) {
			blackCanCastleKingSide = false
			blackCanCastleQueenSide = false
		} else if move.Start.Equals(Position{X: 0, Y: 7}) {
			blackCanCastleQueenSide = false
		} else if move.Start.Equals(Position{X: 7, Y: 7}) {
			blackCanCastleKingSide = false
		}
	} else {
		if move.Start.Equals(Position{X: 4, Y: 0}) {
			whiteCanCastleKingSide = false
			whiteCanCastleQueenSide = false
		} else if move.Start.Equals(Position{X: 0, Y: 0}) {
			whiteCanCastleQueenSide = false
		} else if move.Start.Equals(Position{X: 7, Y: 0}) {
			whiteCanCastleKingSide = false
		}
	}

	return State{
		state.Board.ExecuteMove(move),
		whiteCanCastleKingSide,
		whiteCanCastleQueenSide,
		blackCanCastleKingSide,
		blackCanCastleQueenSide,
		!state.Turn,
		move.EncodeMove(),
	}
}
