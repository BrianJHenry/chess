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
		baseString, err = move.getMoveAlgebraicNotation(state.Board)
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

func (move Move) getMoveAlgebraicNotation(board Board) (string, error) {
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

		// Disambiguate and get core
		visiblePositions := getDirectionalVision(board, move.End, rookDirections)
		core, err := getCoreAlgebraicNotation(board, piece, move, visiblePositions[:])
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteKnight, BlackKnight:
		base := "N"

		// Disambiguate and get core
		visiblePositions := getKnightVision(move.End)
		core, err := getCoreAlgebraicNotation(board, piece, move, ConvertToNullablePositions(visiblePositions))
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteBishop, BlackBishop:
		base := "B"

		// Disambiguate and get core
		visiblePositions := getDirectionalVision(board, move.End, bishopDirections)
		core, err := getCoreAlgebraicNotation(board, piece, move, visiblePositions[:])
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteQueen, BlackQueen:
		base := "Q"

		// Disambiguate and get core
		bishopVisions := getDirectionalVision(board, move.End, bishopDirections)
		rookVisions := getDirectionalVision(board, move.End, rookDirections)
		visiblePositions := bishopVisions[:]
		visiblePositions = append(visiblePositions, rookVisions[:]...)
		core, err := getCoreAlgebraicNotation(board, piece, move, visiblePositions)
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteKing, BlackKing:
		base := "K"

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
	}

	return "", nil
}

func getCoreAlgebraicNotation(board Board, piece Piece, move Move, visiblePositions []NullablePosition) (string, error) {
	base := ""

	// Disambiguate
	sameFile := false
	sameRank := false
	anyAmbiguities := false
	for _, visiblePosition := range visiblePositions {
		// Check for rooks that could make the same move
		if visiblePosition.Valid && board.GetSquare(visiblePosition.Position) == piece {
			if visiblePosition.Position.X == move.Start.X {
				sameRank = true
			} else if visiblePosition.Position.Y == move.Start.Y {
				sameFile = true
			}
			anyAmbiguities = true
		}
	}
	if sameFile && sameRank {
		startPosition, err := ConvertPositionToString(move.Start)
		if err != nil {
			return "", err
		}
		return startPosition, nil
	} else if sameFile {
		base += ConvertRankToString(move.Start.X)
	} else if anyAmbiguities {
		base += ConvertFileToString(move.Start.Y)
	}

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
}

func AlgebraicNotationToMove(algebraicNotation string, state State) (move Move) {
	return
}
