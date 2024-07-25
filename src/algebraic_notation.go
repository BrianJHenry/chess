package chess

func (move Move) ToAlgebraicNotation(state State) (string, error) {
	var baseString string
	switch move.Flag {
	case KingSideCastle:
		baseString = "0-0"
	case QueenSideCastle:
		baseString = "0-0-0"
	default:
		var err error
		baseString, err = move.baseAlgebraicNotation(state.Board)
		if err != nil {
			return "", err
		}
	}

	// Check and mate
	updatedState := state.ExecuteMove(move)
	isInCheck, err := updatedState.Board.IsInCheck(updatedState.Turn)
	if err != nil {
		return baseString, err
	}

	if isInCheck {
		moves, err := GenerateAllMoves(updatedState)
		if err != nil {
			return baseString, err
		}

		if len(moves) == 0 {
			baseString += "#"
		} else {
			baseString += "+"
		}
	}

	return baseString, nil
}

func (move Move) baseAlgebraicNotation(board Board) (string, error) {
	piece := board.GetSquare(move.Start)

	switch piece {
	case WhitePawn, BlackPawn:
		// Base movement
		base, err := ConvertPositionToString(move.End)
		if err != nil {
			return "", err
		}

		// Capture
		if move.Start.Y != move.End.Y {
			base = string(ConvertFileToString(move.Start.Y)) + base
		}

		// Promotion
		switch move.Flag {
		case PromoteToQueen:
			base += "=Q"
		case PromoteToRook:
			base += "=R"
		case PromoteToBishop:
			base += "=B"
		case PromoteToKnight:
			base += "=N"
		}

		return base, nil
	case WhiteRook, BlackRook:
		base := "R"

		// Disambiguate
		visiblePositions := getDirectionalVision(board, move.End, rookDirections)
		disambiguated, err := disambiguateAlgebraicNotation(board, piece, move.Start, visiblePositions[:])
		if err != nil {
			return "", nil
		}
		base += disambiguated

		// Captures
		if board.GetSquare(move.End) != EmptySquare {
			base += "x"
		}

		// End position
		endPosition, err := ConvertPositionToString(move.End)
		if err != nil {
			return "", err
		}

		base += endPosition
		return base, nil
	case WhiteKnight, BlackKnight:
		//
	case WhiteBishop, BlackBishop:
		//
	case WhiteQueen, BlackQueen:
		//
	case WhiteKing, BlackKing:
		//
	}

	return "", nil
}

func disambiguateAlgebraicNotation(board Board, piece Piece, start Position, visiblePositions []Position) (string, error) {
	sameFile := false
	sameRank := false
	anyAmbiguities := false
	for _, visiblePosition := range visiblePositions {
		// Check for rooks that could make the same move
		if visiblePosition != start && board.GetSquare(visiblePosition) == piece {
			if visiblePosition.X == start.X {
				sameRank = true
			} else if visiblePosition.Y == start.Y {
				sameFile = true
			}
			anyAmbiguities = true
		}
	}
	if sameFile && sameRank {
		startPosition, err := ConvertPositionToString(start)
		if err != nil {
			return "", err
		}
		return startPosition, nil
	} else if sameFile {
		return ConvertRankToString(start.X), nil
	} else if anyAmbiguities {
		return ConvertFileToString(start.Y), nil
	}

	return "", nil
}

func AlgebraicNotationToMove(algebraicNotation string, board Board) (move Move) {
	return
}
