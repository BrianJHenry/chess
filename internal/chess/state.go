package chess

type Color bool

const (
	White Color = false
	Black Color = true
)

type State struct {
	Board             Board
	CastlingRights    CastlingRights
	ActiveColor       Color
	EnPassantPosition PositionOpt
}

type CastlingRights struct {
	WhiteCanCastleKingSide  bool
	WhiteCanCastleQueenSide bool
	BlackCanCastleKingSide  bool
	BlackCanCastleQueenSide bool
}

func InitialiseState() State {
	return State{
		InitialiseBoard(),
		CastlingRights{
			true,
			true,
			true,
			true,
		},
		White,
		PositionOpt{Ok: false},
	}
}

// DoMove takes in a state and a move and executes the move, returning the updated state.
func (state *State) DoMove(move Move) {
	castlingRights := state.CastlingRights
	whiteCanCastleKingSide := castlingRights.WhiteCanCastleKingSide
	whiteCanCastleQueenSide := castlingRights.WhiteCanCastleQueenSide
	blackCanCastleKingSide := castlingRights.BlackCanCastleKingSide
	blackCanCastleQueenSide := castlingRights.BlackCanCastleQueenSide

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

	state.Board.DoMove(move)
	state.CastlingRights = CastlingRights{
		whiteCanCastleKingSide,
		whiteCanCastleQueenSide,
		blackCanCastleKingSide,
		blackCanCastleQueenSide,
	}
	state.ActiveColor = !state.ActiveColor
	state.EnPassantPosition = enPassantSquare
}

func (state *State) UndoMove(move Move, castlingRights CastlingRights, enPassantSquare PositionOpt) {
	state.Board.UndoMove(move)
	state.CastlingRights = castlingRights
	state.ActiveColor = !state.ActiveColor
	state.EnPassantPosition = enPassantSquare
}
