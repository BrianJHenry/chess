package chess

import (
	"errors"
	"fmt"
	"strconv"
)

// AlgebraicNotation represents the standard written notation for a chess move.
type AlgebraicNotation string

// ToAlgebraicNotation converts from a move struct to standard algebraic notation.
func (move Move) ToAlgebraicNotation(state State) (AlgebraicNotation, error) {
	var baseString string
	switch move.Flag {
	case KingSideCastle:
		baseString = "O-O"
	case QueenSideCastle:
		baseString = "O-O-O"
	default:
		var err error
		baseString, err = move.getAlgebraicNotationCore(state.Board)
		if err != nil {
			return "", err
		}
	}

	// Check and mate
	updatedState := state.DoMove(move)
	isInCheck, err := isInCheck(updatedState.Board, updatedState.ActiveColor)
	if err != nil {
		return AlgebraicNotation(baseString), err
	}
	if isInCheck {
		moves, err := GenerateAllMoves(updatedState)
		if err != nil {
			return AlgebraicNotation(baseString), err
		}

		if len(moves) == 0 {
			baseString += "#"
		} else {
			baseString += "+"
		}
	}

	return AlgebraicNotation(baseString), nil
}

// ToMove converts from standard algebraic notation to a move struct.
func (algebraicNotation AlgebraicNotation) ToMove(state State) (Move, error) {
	// Castling
	if algebraicNotation == "O-O" || algebraicNotation == "O-O+" || algebraicNotation == "O-O#" {
		var rank int8
		if state.ActiveColor == Black {
			rank = 0
		} else {
			rank = 7
		}
		return Move{
			Start: Position{rank, 4},
			End:   Position{rank, 6},
			Flag:  KingSideCastle,
		}, nil
	} else if algebraicNotation == "O-O-O" || algebraicNotation == "O-O-O+" || algebraicNotation == "O-O-O#" {
		var rank int8
		if state.ActiveColor == Black {
			rank = 0
		} else {
			rank = 7
		}
		return Move{
			Start: Position{rank, 4},
			End:   Position{rank, 2},
			Flag:  QueenSideCastle,
		}, nil
	}

	length := len(algebraicNotation)
	lastChar := algebraicNotation[length-1]
	if lastChar == '#' || lastChar == '+' {
		length--
	}

	var flag MoveFlag = None
	if algebraicNotation[length-2] == '=' {
		switch algebraicNotation[length-1] {
		case 'Q':
			flag = PromoteToQueen
		case 'R':
			flag = PromoteToRook
		case 'B':
			flag = PromoteToBishop
		case 'N':
			flag = PromoteToKnight
		}
		length -= 2
	}

	an := string(algebraicNotation)
	endPosition, err := stringToPosition(an[length-2 : length])
	if err != nil {
		return Move{}, err
	}
	length -= 2

	if length == 0 {
		piece := getPiece('P', state.ActiveColor)
		start, err := findStartPosition(piece, state, endPosition, -1, -1)
		if err != nil {
			return Move{}, err
		}
		return Move{
			start,
			endPosition,
			None,
		}, nil
	}

	isCapture := false
	if algebraicNotation[length-1] == 'x' {
		isCapture = true
		length--
	}

	piece := getPiece(algebraicNotation[0], state.ActiveColor)
	var start Position
	if length == 1 && piece != WhitePawn && piece != BlackPawn {
		start, err = findStartPosition(piece, state, endPosition, -1, -1)
		if err != nil {
			return Move{}, err
		}
	} else if length == 1 {
		file, err := runeToFile(rune(algebraicNotation[0]))
		if err != nil {
			return Move{}, err
		}

		start, err = findStartPosition(piece, state, endPosition, -1, file)
		if err != nil {
			return Move{}, err
		}
	} else if length == 2 {
		if rank, err := strconv.Atoi(string(algebraicNotation[1])); err == nil {
			start, err = findStartPosition(piece, state, endPosition, int8(rank), -1)
			if err != nil {
				return Move{}, err
			}
		} else {
			file, err := runeToFile(rune(algebraicNotation[1]))
			if err != nil {
				return Move{}, err
			}
			start, err = findStartPosition(piece, state, endPosition, -1, file)
			if err != nil {
				return Move{}, err
			}
		}
	} else {
		// Case where the disambiguating gives the full position
		start, err = stringToPosition(string(algebraicNotation)[1:2])
		if err != nil {
			return Move{}, err
		}
	}

	// Check for en passant
	if isCapture && (piece == WhitePawn || piece == BlackPawn) && state.Board.GetSquare(endPosition) == EmptySquare {
		flag = EnPassant
	}

	return Move{
		Start: start,
		End:   endPosition,
		Flag:  flag,
	}, nil
}

// getAlgebraicNotationCore generates the core notation for non-castling moves.
func (move Move) getAlgebraicNotationCore(board Board) (string, error) {
	piece := board.GetSquare(move.Start)

	switch piece {
	case WhitePawn, BlackPawn:
		// Base movement
		base, err := positionToString(move.End)
		if err != nil {
			return "", err
		}

		// Capture
		if move.Start.Y != move.End.Y {
			file, err := fileToString(move.Start.Y)
			if err != nil {
				return "", err
			}
			base = file + "x" + base
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
		core, err := getPositionalAlgebraicNotation(board, piece, move, visiblePositions[:])
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteKnight, BlackKnight:
		base := "N"

		// Disambiguate and get core
		visiblePositions := getKnightVision(move.End)
		core, err := getPositionalAlgebraicNotation(board, piece, move, PositionsToOptionalPositions(visiblePositions))
		if err != nil {
			return "", err
		}
		base += core

		return base, nil
	case WhiteBishop, BlackBishop:
		base := "B"

		// Disambiguate and get core
		visiblePositions := getDirectionalVision(board, move.End, bishopDirections)
		core, err := getPositionalAlgebraicNotation(board, piece, move, visiblePositions[:])
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
		core, err := getPositionalAlgebraicNotation(board, piece, move, visiblePositions)
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
		endPosition, err := positionToString(move.End)
		if err != nil {
			return "", err
		}
		base += endPosition

		return base, nil
	}

	return "", nil
}

// getPositionalAlgebraicNotation generates the positional element of algebraic notation: i.e. xa2 or b3xc4.
func getPositionalAlgebraicNotation(board Board, piece Piece, move Move, visiblePositions []PositionOpt) (string, error) {
	base := ""

	// Disambiguate
	sameFile := false
	sameRank := false
	anyAmbiguities := false
	for _, visiblePosition := range visiblePositions {
		// Check for rooks that could make the same move
		if visiblePosition.Ok && visiblePosition.Position != move.Start && board.GetSquare(visiblePosition.Position) == piece {
			if visiblePosition.Position.X == move.Start.X {
				sameRank = true
			} else if visiblePosition.Position.Y == move.Start.Y {
				sameFile = true
			}
			anyAmbiguities = true
		}
	}
	if sameFile && sameRank {
		startPosition, err := positionToString(move.Start)
		if err != nil {
			return "", err
		}
		return startPosition, nil
	} else if sameFile {
		rank, err := rankToString(move.Start.X)
		if err != nil {
			return "", err
		}
		base += rank
	} else if anyAmbiguities {
		file, err := fileToString(move.Start.Y)
		if err != nil {
			return "", err
		}
		base += file
	}

	// Captures
	if board.GetSquare(move.End) != EmptySquare {
		base += "x"
	}

	// End position
	endPosition, err := positionToString(move.End)
	if err != nil {
		return "", err
	}
	base += endPosition

	return base, nil
}

// findStartPosition parses move.Start from a move in algebraic notation.
func findStartPosition(piece Piece, state State, end Position, hintX, hintY int8) (Position, error) {
	switch piece {
	case WhitePawn, BlackPawn:
		var searchDirection int8
		if piece == WhitePawn {
			searchDirection = 1
		} else {
			searchDirection = -1
		}

		// Capture
		if hintY != -1 {
			return Position{end.X + searchDirection, hintY}, nil
		}

		// Non-Capture
		oneSquare := Position{end.X + searchDirection, end.Y}
		if state.Board.GetSquare(oneSquare) == piece {
			return oneSquare, nil
		} else {
			return Position{end.X + (2 * searchDirection), end.Y}, nil
		}
	case WhiteRook, BlackRook:
		rookVisions := getDirectionalVision(state.Board, end, rookDirections)
		start, err := findStartPositionForVisions(piece, state.Board, hintX, hintY, rookVisions[:])
		if err != nil {
			return Position{}, err
		}
		return start, nil
	case WhiteKnight, BlackKnight:
		knighVisions := getKnightVision(end)
		start, err := findStartPositionForVisions(piece, state.Board, hintX, hintY, PositionsToOptionalPositions(knighVisions))
		if err != nil {
			return Position{}, err
		}
		return start, nil
	case WhiteBishop, BlackBishop:
		bishopVisions := getDirectionalVision(state.Board, end, bishopDirections)
		start, err := findStartPositionForVisions(piece, state.Board, hintX, hintY, bishopVisions[:])
		if err != nil {
			return Position{}, err
		}
		return start, nil
	case WhiteQueen, BlackQueen:
		rookVisions := getDirectionalVision(state.Board, end, rookDirections)
		bishopVisions := getDirectionalVision(state.Board, end, bishopDirections)
		visions := rookVisions[:]
		visions = append(visions, bishopVisions[:]...)
		start, err := findStartPositionForVisions(piece, state.Board, hintX, hintY, visions)
		if err != nil {
			return Position{}, err
		}
		return start, nil
	case WhiteKing, BlackKing:
		start, err := findKing(state.Board, state.ActiveColor)
		if err != nil {
			return Position{}, err
		}
		return start, nil
	}
	return Position{}, errors.New("invalid piece")
}

// findStartPositionForVisions retrieves move.Start based on the possible moves of the moved piece type.
func findStartPositionForVisions(piece Piece, board Board, hintX, hintY int8, visions []PositionOpt) (Position, error) {
	for _, vision := range visions {
		// Skip out of bounds
		if !vision.Ok {
			continue
		}

		// TODO: improve this ugliness!
		if hintX != -1 {
			if vision.Position.X == hintX && board.GetSquare(vision.Position) == piece {
				return vision.Position, nil
			}
		} else if hintY != -1 {
			if vision.Position.Y == hintY && board.GetSquare(vision.Position) == piece {
				return vision.Position, nil
			}
		} else if board.GetSquare(vision.Position) == piece {
			return vision.Position, nil
		}
	}

	return Position{}, errors.New("given the constraints, the starting position could not be found")
}

// getPiece retrieves the Piece type based on the turn and first character in the algebraic notation version of the move.
func getPiece(piece byte, turn Color) Piece {
	switch piece {
	case 'K':
		if turn == Black {
			return BlackKing
		} else {
			return WhiteKing
		}
	case 'Q':
		if turn == Black {
			return BlackQueen
		} else {
			return WhiteQueen
		}
	case 'R':
		if turn == Black {
			return BlackRook
		} else {
			return WhiteRook
		}
	case 'B':
		if turn == Black {
			return BlackBishop
		} else {
			return WhiteBishop
		}
	case 'N':
		if turn == Black {
			return BlackKnight
		} else {
			return WhiteKnight
		}
	default: // pawn
		if turn == Black {
			return BlackPawn
		} else {
			return WhitePawn
		}
	}
}

// stringToPosition takes in a 2 character letter number combo to specify the file and rank of a position. Ex. A1 -> 7, 0; C4 => 4, 2.
func stringToPosition(pos string) (Position, error) {
	if len(pos) != 2 {
		return Position{}, errors.New("position codes should be of length 2")
	}

	rankChar := pos[1]
	rank := 7 - (rankChar - '1')

	fileChar := pos[0]
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

// positionToString takes in a position and returns the 2 character letter number combo that specifies the file and rank of the position.
func positionToString(position Position) (string, error) {
	if !isInBounds(position) {
		return "", fmt.Errorf("position out of bounds: %d, %d", position.X, position.Y)
	}

	rankChar := (7 - position.X) + '1'
	fileChar := position.Y + 'a'

	return fmt.Sprintf("%c%c", fileChar, rankChar), nil
}

// runeToFile takes in a rune from a-h and returns the index of the file.
func runeToFile(f rune) (int8, error) {
	file := int8(f - 'a')
	if file < 0 || file > 7 {
		return file, fmt.Errorf("file %c is out of range", f)
	}
	return file, nil
}

// fileToString takes in the index of a file and returns the character from a-h that represents the file.
func fileToString(f int8) (string, error) {
	if f < 0 || f > 7 {
		return "", fmt.Errorf("file %d is out of range", f)
	}
	return string(rune(f + 'a')), nil
}

// rankToString takes in the index of the rank and returns the character from 1-8 that represents the rank.
func rankToString(r int8) (string, error) {
	if r < 0 || r > 7 {
		return "", fmt.Errorf("rank %d is out of range", r)
	}
	return fmt.Sprint(7 - r + 1), nil
}
