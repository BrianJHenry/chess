package chess

type Piece int8

const (
	WhiteKing Piece = iota - 6
	WhiteQueen
	WhiteRook
	WhiteBishop
	WhiteKnight
	WhitePawn
	EmptySquare
	BlackPawn
	BlackKnight
	BlackBishop
	BlackRook
	BlackQueen
	BlackKing
)

func getRookColorForKing(king Piece) Piece {
	if king == BlackKing {
		return BlackRook
	} else {
		return WhiteRook
	}
}

func getQueenColorForPawn(pawn Piece) Piece {
	if pawn == BlackPawn {
		return BlackQueen
	} else {
		return WhiteQueen
	}
}

func getRookColorForPawn(pawn Piece) Piece {
	if pawn == BlackPawn {
		return BlackRook
	} else {
		return WhiteRook
	}
}

func getBishopColorForPawn(pawn Piece) Piece {
	if pawn == BlackPawn {
		return BlackBishop
	} else {
		return WhiteBishop
	}
}

func getKnightColorForPawn(pawn Piece) Piece {
	if pawn == BlackPawn {
		return BlackKnight
	} else {
		return WhiteKnight
	}
}
