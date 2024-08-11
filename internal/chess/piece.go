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

// PieceToDisplayString returns a 2 character string representation of a piece to be displayed as ascii.
func PieceToDisplayString(piece Piece) string {
	switch piece {
	case WhitePawn:
		return "WP"
	case WhiteKnight:
		return "WN"
	case WhiteBishop:
		return "WB"
	case WhiteRook:
		return "WR"
	case WhiteQueen:
		return "WQ"
	case WhiteKing:
		return "WK"
	case BlackPawn:
		return "BP"
	case BlackKnight:
		return "BN"
	case BlackBishop:
		return "BB"
	case BlackRook:
		return "BR"
	case BlackQueen:
		return "BQ"
	case BlackKing:
		return "BK"
	default:
		return "  "
	}
}

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

func getEnemyPawnColor(color Color) Piece {
	if color == Black {
		return WhitePawn
	} else {
		return BlackPawn
	}
}

func getSameColorPawn(piece Piece) Piece {
	if piece < 0 {
		return WhitePawn
	} else if piece > 0 {
		return BlackPawn
	} else {
		return EmptySquare
	}
}
