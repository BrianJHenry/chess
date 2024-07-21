package chess

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

func (position Position) AddOffset(offset Position) Position {
	return Position{X: position.X + offset.X, Y: position.Y + offset.Y}
}

func (position Position) MultiplyScalar(scalar int8) Position {
	return Position{X: position.X * scalar, Y: position.Y * scalar}
}
