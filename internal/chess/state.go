package chess

type Color bool

const (
	White Color = false
	Black Color = true
)

type State struct {
	Board                   Board
	WhiteCanCastleKingSide  bool
	WhiteCanCastleQueenSide bool
	BlackCanCastleKingSide  bool
	BlackCanCastleQueenSide bool
	ActiveColor             Color
	EnPassantPosition       PositionOpt
}

func InitialiseState() State {
	return State{
		InitialiseBoard(),
		true,
		true,
		true,
		true,
		White,
		PositionOpt{Ok: false},
	}
}

// DoMove takes in a state and a move and executes the move, returning the updated state.
func (state State) DoMove(move Move) State {
	whiteCanCastleKingSide := state.WhiteCanCastleKingSide
	whiteCanCastleQueenSide := state.WhiteCanCastleQueenSide
	blackCanCastleKingSide := state.BlackCanCastleKingSide
	blackCanCastleQueenSide := state.BlackCanCastleQueenSide

	if state.ActiveColor == Black && (move.Start == Position{0, 4}) {
		blackCanCastleKingSide = false
		blackCanCastleQueenSide = false
	} else if state.ActiveColor == White && (move.Start == Position{7, 4}) {
		whiteCanCastleKingSide = false
		whiteCanCastleQueenSide = false
	}

	// Check for rooks moving or being captured
	if MoveTouchesSquare(move, Position{0, 0}) {
		blackCanCastleQueenSide = false
	}
	if MoveTouchesSquare(move, Position{0, 7}) {
		blackCanCastleKingSide = false
	}
	if MoveTouchesSquare(move, Position{7, 0}) {
		whiteCanCastleQueenSide = false
	}
	if MoveTouchesSquare(move, Position{7, 7}) {
		whiteCanCastleKingSide = false
	}

	enPassantSquare := PositionOpt{
		Ok: false,
	}
	if (state.ActiveColor == Black && state.Board.GetSquare(move.Start) == BlackPawn && (move.End.X-move.Start.X == 2)) ||
		(state.ActiveColor == White && state.Board.GetSquare(move.Start) == WhitePawn && (move.Start.X-move.End.X == 2)) {

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

func (state State) UndoMove(move Move, whiteCastleQ bool, whiteCastleK bool, blackCastleQ bool, blackCastleK bool, enPassantSquare PositionOpt) State {
	return State{
		state.Board.UndoMove(move),
		whiteCastleK,
		whiteCastleQ,
		blackCastleK,
		blackCastleQ,
		!state.ActiveColor,
		enPassantSquare,
	}
}
