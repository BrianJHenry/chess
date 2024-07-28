package chess

type ActiveColor bool

const (
	WhiteTurn ActiveColor = false
	BlackTurn ActiveColor = true
)

type State struct {
	Board                   Board
	WhiteCanCastleKingSide  bool
	WhiteCanCastleQueenSide bool
	BlackCanCastleKingSide  bool
	BlackCanCastleQueenSide bool
	ActiveColor             ActiveColor
	EnPassantPosition       PositionOpt
}

// DoMove takes in a state and a move and executes the move, returning the updated state.
func (state State) DoMove(move Move) State {
	whiteCanCastleKingSide := state.WhiteCanCastleKingSide
	whiteCanCastleQueenSide := state.WhiteCanCastleQueenSide
	blackCanCastleKingSide := state.BlackCanCastleKingSide
	blackCanCastleQueenSide := state.BlackCanCastleQueenSide

	if state.ActiveColor == BlackTurn {
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

	enPassantSquare := PositionOpt{
		Ok: false,
	}
	if (state.ActiveColor == BlackTurn && state.Board.GetSquare(move.Start) == BlackPawn && (move.End.X-move.Start.X == 2)) ||
		(state.ActiveColor == WhiteTurn && state.Board.GetSquare(move.Start) == WhitePawn && (move.Start.X-move.End.X == 2)) {

		enPassantSquare.Ok = true
		enPassantSquare.Position = Position{
			X: (move.Start.X + move.End.X) / 2,
			Y: move.End.Y,
		}
	}

	return State{
		state.Board.DoMove(move),
		whiteCanCastleKingSide,
		whiteCanCastleQueenSide,
		blackCanCastleKingSide,
		blackCanCastleQueenSide,
		!state.ActiveColor,
		enPassantSquare,
	}
}
