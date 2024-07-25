package chess

import (
	"errors"
	"fmt"
)

type Board [8][8]Piece

type Position struct {
	X, Y int8
}

type NullablePosition struct {
	Position Position
	Valid    bool
}

func (board Board) GetPrintableBoard() string {
	stringBoard := ""
	for i := 0; i < 8; i++ {
		stringBoard += "+----+----+----+----+----+----+----+----+\n"
		stringBoard += "|    |    |    |    |    |    |    |    |\n"
		for j := 0; j < 8; j++ {
			stringBoard += fmt.Sprintf("| %s ", board[i][j].GetPrintablePiece())
		}
		stringBoard += "|\n"
		stringBoard += "|    |    |    |    |    |    |    |    |\n"
		stringBoard += "|    |    |    |    |    |    |    |    |\n"
	}
	stringBoard += "+----+----+----+----+----+----+----+----+\n"
	return stringBoard
}

func InitialPosition() Board {
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

// Update board for move
func (board Board) ExecuteMove(move Move) Board {
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

func (board Board) GetSquare(position Position) Piece {
	return board[position.X][position.Y]
}

func (board Board) FindKing(color Turn) (kingPosition Position, err error) {
	var square Piece

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			kingPosition = Position{i, j}
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
	return Position{position.X + offset.X, position.Y + offset.Y}
}

func (position Position) MultiplyScalar(scalar int8) Position {
	return Position{position.X * scalar, position.Y * scalar}
}

// Ex. A1 -> 7, 0; C4 => 4, 2
func ConvertStringToPosition(stringPosition string) (Position, error) {
	if len(stringPosition) != 2 {
		return Position{}, errors.New("position codes should be of length 2")
	} else {
		rankChar := stringPosition[1]
		rank := 7 - (rankChar - '1')

		fileChar := stringPosition[0]
		file := fileChar - 'a'
		position := Position{
			int8(rank),
			int8(file),
		}

		if isInBounds(position) {
			return position, nil
		} else {
			return position, fmt.Errorf("invalid position: %d, %d", position.X, position.Y)
		}
	}
}
