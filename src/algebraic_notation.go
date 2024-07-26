package chess

import (
	"errors"
	"fmt"
	"strconv"
)

func (move Move) ToAlgebraicNotation(state State) (string, error) {
	var baseString string
	switch move.Flag {
	case KingSideCastle:
		baseString = "O-O"
	case QueenSideCastle:
		baseString = "O-O-O"
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
			base = string(ConvertFileToString(move.Start.Y)) + "x" + base
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
		if visiblePosition.Valid && visiblePosition.Position != move.Start && board.GetSquare(visiblePosition.Position) == piece {
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

func AlgebraicNotationToMove(algebraicNotation string, state State) (Move, error) {
	// Castling
	if algebraicNotation == "O-O" || algebraicNotation == "O-O+" || algebraicNotation == "O-O#" {
		var rank int8
		if state.Turn == BlackTurn {
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
		if state.Turn == BlackTurn {
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

	endPosition, err := ConvertStringToPosition(algebraicNotation[length-2 : length-1])
	if err != nil {
		return Move{}, err
	}
	length -= 2

	isCapture := false
	if algebraicNotation[length-1] == 'x' {
		isCapture = true
		length--
	}

	var start Position
	if length == 1 {
		start, err = findStartPosition(algebraicNotation[0], state, endPosition, -1, -1)
		if err != nil {
			return Move{}, err
		}
	} else if length == 2 {
		if rank, err := strconv.Atoi(string(algebraicNotation[1])); err == nil {
			start, err = findStartPosition(algebraicNotation[0], state, endPosition, int8(rank), -1)
			if err != nil {
				return Move{}, err
			}
		} else {
			start, err = findStartPosition(algebraicNotation[0], state, endPosition, -1, ConvertRuneToFile(rune(algebraicNotation[1])))
			if err != nil {
				return Move{}, err
			}
		}
	} else {
		// Case where the disambiguating gives the full position
		start, err = ConvertStringToPosition(algebraicNotation[1:2])
		if err != nil {
			return Move{}, err
		}
	}

	// Check for en passant
	piece := state.Board.GetSquare(start)
	if isCapture && piece == WhitePawn || piece == BlackPawn && state.Board.GetSquare(endPosition) == EmptySquare {
		flag = EnPassant
	}

	return Move{
		Start: start,
		End:   endPosition,
		Flag:  flag,
	}, nil
}

func findStartPosition(charPiece byte, state State, end Position, hintX, hintY int8) (Position, error) {
	piece := getPiece(charPiece, state.Turn)
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
		start, err := findStartPositionForVisions(piece, state.Board, hintX, hintY, ConvertToNullablePositions(knighVisions))
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
		start, err := state.Board.FindKing(state.Turn)
		if err != nil {
			return Position{}, err
		}
		return start, nil
	}
	return Position{}, fmt.Errorf("invalid piece: %c", rune(charPiece))
}

func findStartPositionForVisions(piece Piece, board Board, hintX, hintY int8, visions []NullablePosition) (Position, error) {
	for _, vision := range visions {
		// Skip out of bounds
		if !vision.Valid {
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

func getPiece(piece byte, turn Turn) Piece {
	switch piece {
	case 'K':
		if turn == BlackTurn {
			return BlackKing
		} else {
			return WhiteKing
		}
	case 'Q':
		if turn == BlackTurn {
			return BlackQueen
		} else {
			return WhiteQueen
		}
	case 'R':
		if turn == BlackTurn {
			return BlackRook
		} else {
			return WhiteRook
		}
	case 'B':
		if turn == BlackTurn {
			return BlackBishop
		} else {
			return WhiteBishop
		}
	case 'N':
		if turn == BlackTurn {
			return BlackKnight
		} else {
			return WhiteKnight
		}
	default: // pawn
		if turn == BlackTurn {
			return BlackPawn
		} else {
			return WhitePawn
		}
	}
}
