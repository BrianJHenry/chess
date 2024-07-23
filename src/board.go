package chess

import "errors"

type Board [8][8]Piece

type Position struct {
	X, Y int8
}

func InitialPosition() Board {
	return Board{
		{WhiteRook, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackRook},
		{WhiteKnight, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackKnight},
		{WhiteBishop, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackBishop},
		{WhiteQueen, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackQueen},
		{WhiteKing, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackKing},
		{WhiteBishop, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackBishop},
		{WhiteKnight, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackKnight},
		{WhiteRook, WhitePawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, BlackPawn, BlackRook},
	}
}

// Update board for move
func (board Board) ExecuteMove(move Move) Board {
	piece := board[move.Start.X][move.Start.Y]

	switch move.Flag {
	case None:
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = piece
	case EnPassant:
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = piece
	case QueenSideCastle:
		queenSideRook := getRookColorForKing(piece)
		board[4][move.Start.Y] = EmptySquare
		board[3][move.Start.Y] = queenSideRook
		board[2][move.Start.Y] = piece
		board[0][move.Start.Y] = EmptySquare
	case KingSideCastle:
		kingSideRook := getRookColorForKing(piece)
		board[4][move.Start.Y] = EmptySquare
		board[5][move.Start.Y] = kingSideRook
		board[6][move.Start.Y] = piece
		board[7][move.Start.Y] = EmptySquare
	case PromoteToQueen:
		queen := getQueenColorForPawn(piece)
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = queen
	case PromoteToRook:
		rook := getRookColorForPawn(piece)
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = rook
	case PromoteToBishop:
		bishop := getBishopColorForPawn(piece)
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = bishop
	case PromoteToKnight:
		knight := getKnightColorForPawn(piece)
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = knight
	}

	return board
}

func (board Board) GetSquare(position Position) Piece {
	return board[position.X][position.Y]
}

func (board Board) FindKing(color Turn) (kingPosition Position, err error) {
	var square Piece

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			kingPosition = Position{X: i, Y: j}
			square = board.GetSquare(kingPosition)

			if (color == BlackTurn && square == BlackKing) ||
				color == WhiteTurn && square == WhiteKing {

				return kingPosition, nil
			}
		}
	}

	return Position{}, errors.New("missing king")
}

func (board Board) IsInCheck(color Turn) (bool, error) {
	kingPosition, err := board.FindKing(color)
	if err != nil {
		return false, err
	}

	return IsSquareAttacked(board, kingPosition, color), nil
}

func (position Position) AddOffset(offset Position) Position {
	return Position{X: position.X + offset.X, Y: position.Y + offset.Y}
}

func (position Position) MultiplyScalar(scalar int8) Position {
	return Position{X: position.X * scalar, Y: position.Y * scalar}
}

func (position Position) Equals(other Position) bool {
	return position.X == other.X && position.Y == other.Y
}
