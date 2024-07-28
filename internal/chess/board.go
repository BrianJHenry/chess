package chess

import (
	"fmt"
)

// Board represents a chess board.
type Board [8][8]Piece

// Position represents a location on a Board.
type Position struct {
	X, Y int8
}

// PositionOpt represents a Position that may or may not be valid.
type PositionOpt struct {
	Position Position
	Ok       bool
}

// InitialiseBoard returns the default starting position of a chess board.
func InitialiseBoard() Board {
	return Board{
		{BlackRook, BlackKnight, BlackBishop, BlackQueen, BlackKing, BlackBishop, BlackKnight, BlackRook},
		{BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn, BlackPawn},
		{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
		{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
		{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
		{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
		{WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn, WhitePawn},
		{WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, WhiteBishop, WhiteKnight, WhiteRook},
	}
}

// PositionsToOptionalPositions converts a slice of Position ot a slice of PositionOpt where Ok is true.
func PositionsToOptionalPositions(positions []Position) []PositionOpt {
	optPositions := make([]PositionOpt, len(positions))

	for i, pos := range positions {
		optPositions[i] = PositionOpt{
			Ok:       true,
			Position: pos,
		}
	}

	return optPositions
}

// BoardToDisplayString converts a position on the board to an visual representation in ascii.
func BoardToDisplayString(board Board) string {
	stringBoard := ""
	for i := 0; i < 8; i++ {
		stringBoard += "  +----+----+----+----+----+----+----+----+\n"
		stringBoard += "  |    |    |    |    |    |    |    |    |\n"
		stringBoard += fmt.Sprintf("%d ", 8-i)
		for j := 0; j < 8; j++ {
			stringBoard += fmt.Sprintf("| %s ", PieceToDisplayString(board[i][j]))
		}
		stringBoard += "|\n"
		stringBoard += "  |    |    |    |    |    |    |    |    |\n"
	}
	stringBoard += "  +----+----+----+----+----+----+----+----+\n"
	stringBoard += "    a    b    c    d    e    f    g    h"
	return stringBoard
}

// DoMove takes in a board and a move and executes the move, returning the updated board.
func (board Board) DoMove(move Move) Board {
	piece := board[move.Start.X][move.Start.Y]

	switch move.Flag {
	case None:
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.End.X][move.End.Y] = piece
	case EnPassant:
		board[move.Start.X][move.Start.Y] = EmptySquare
		board[move.Start.X][move.End.Y] = EmptySquare
		board[move.End.X][move.End.Y] = piece
	case QueenSideCastle:
		queenSideRook := getRookColorForKing(piece)
		board[move.Start.X][4] = EmptySquare
		board[move.Start.X][3] = queenSideRook
		board[move.Start.X][2] = piece
		board[move.Start.X][0] = EmptySquare
	case KingSideCastle:
		kingSideRook := getRookColorForKing(piece)
		board[move.Start.X][4] = EmptySquare
		board[move.Start.X][5] = kingSideRook
		board[move.Start.X][6] = piece
		board[move.Start.X][7] = EmptySquare
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

// GetSquare returns the piece at the position.
func (board Board) GetSquare(position Position) Piece {
	return board[position.X][position.Y]
}

// AddPositions returns the component-wise addition of two positions.
func AddPositions(position, offset Position) Position {
	return Position{position.X + offset.X, position.Y + offset.Y}
}

// MultiplyScalar multiplies each component of the position by the scalar value.
func MultiplyScalar(position Position, scalar int8) Position {
	return Position{position.X * scalar, position.Y * scalar}
}
