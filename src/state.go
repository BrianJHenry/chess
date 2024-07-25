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
	EnPassantPosition       NullablePosition
}

func (state State) ExecuteMove(move Move) State {
	whiteCanCastleKingSide := state.WhiteCanCastleKingSide
	whiteCanCastleQueenSide := state.WhiteCanCastleQueenSide
	blackCanCastleKingSide := state.BlackCanCastleKingSide
	blackCanCastleQueenSide := state.BlackCanCastleQueenSide

	if state.Turn == BlackTurn {
		if (move.Start == Position{0, 4}) {
			blackCanCastleKingSide = false
			blackCanCastleQueenSide = false
		} else if (move.Start == Position{0, 0}) {
			blackCanCastleQueenSide = false
		} else if (move.Start == Position{0, 7}) {
			blackCanCastleKingSide = false
		}
	} else {
		if (move.Start == Position{7, 4}) {
			whiteCanCastleKingSide = false
			whiteCanCastleQueenSide = false
		} else if (move.Start == Position{7, 0}) {
			whiteCanCastleQueenSide = false
		} else if (move.Start == Position{7, 7}) {
			whiteCanCastleKingSide = false
		}
	}

	enPassantSquare := NullablePosition{
		Valid: false,
	}
	if (state.Turn == BlackTurn && state.Board.GetSquare(move.Start) == BlackPawn && (move.End.X-move.Start.X == 2)) ||
		(state.Turn == WhiteTurn && state.Board.GetSquare(move.Start) == WhitePawn && (move.Start.X-move.End.X == 2)) {

		enPassantSquare.Valid = true
		enPassantSquare.Position = Position{
			X: (move.Start.X + move.End.X) / 2,
			Y: move.End.Y,
		}
	}

	return State{
		state.Board.ExecuteMove(move),
		whiteCanCastleKingSide,
		whiteCanCastleQueenSide,
		blackCanCastleKingSide,
		blackCanCastleQueenSide,
		!state.Turn,
		enPassantSquare,
	}
}

func (state State) Equals(other State) bool {
	return state.Board == other.Board &&
		state.WhiteCanCastleKingSide == other.WhiteCanCastleKingSide &&
		state.BlackCanCastleKingSide == other.BlackCanCastleKingSide &&
		state.WhiteCanCastleQueenSide == other.WhiteCanCastleQueenSide &&
		state.BlackCanCastleQueenSide == other.BlackCanCastleQueenSide &&
		state.Turn == other.Turn
}
