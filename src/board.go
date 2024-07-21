package chess

import "errors"

type Board [8][8]Piece

type Position struct {
	X, Y int8
}

func InitialPosition() *Board {
	return &Board{
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

// TODO
func (board Board) ExecuteMove(move Move) Board {
	return Board{}
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

	return Position{}, errors.New("Missing king")
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
