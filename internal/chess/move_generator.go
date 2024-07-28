package chess

import "errors"

var rookDirections [4]Position = [4]Position{
	{X: 1, Y: 0},
	{X: -1, Y: 0},
	{X: 0, Y: 1},
	{X: 0, Y: -1},
}

var bishopDirections [4]Position = [4]Position{
	{X: 1, Y: 1},
	{X: 1, Y: -1},
	{X: -1, Y: 1},
	{X: -1, Y: -1},
}

var knightOffsets [8]Position = [8]Position{
	{X: 2, Y: 1},
	{X: 1, Y: 2},
	{X: -2, Y: 1},
	{X: -1, Y: 2},
	{X: -2, Y: -1},
	{X: -1, Y: -2},
	{X: 2, Y: -1},
	{X: 1, Y: -2},
}

// GenerateAllMoves generates all possible moves in a given state.
func GenerateAllMoves(state State) (moves []Move, err error) {
	kingPosition, err := findKing(state.Board, state.ActiveColor)
	if err != nil {
		return moves, err
	}

	for i := int8(0); i < 8; i++ {
		for j := int8(0); j < 8; j++ {
			position := Position{i, j}
			square := state.Board.GetSquare(position)

			// Skip any squares that aren't of the active color
			if (state.ActiveColor == Black && square <= 0) || (state.ActiveColor == White && square >= 0) {
				continue
			}

			switch square {
			case WhitePawn, BlackPawn:
				moves = append(moves, GeneratePawnMoves(state, position, kingPosition)...)
			case WhiteBishop, BlackBishop:
				moves = append(moves, GenerateBishopMoves(state, position, kingPosition)...)
			case WhiteKnight, BlackKnight:
				moves = append(moves, GenerateKnightMoves(state, position, kingPosition)...)
			case WhiteRook, BlackRook:
				moves = append(moves, GenerateRookMoves(state, position, kingPosition)...)
			case WhiteQueen, BlackQueen:
				moves = append(moves, GenerateQueenMoves(state, position, kingPosition)...)
			case WhiteKing, BlackKing:
				moves = append(moves, GenerateKingMoves(state, position)...)
			default:
				// Pass on empty squares
			}
		}
	}

	return
}

func GenerateKingMoves(state State, position Position) (moves []Move) {
	piece := state.Board[position.X][position.Y]

	directions := make([]Position, 8)
	directions = append(directions, bishopDirections[:]...)
	directions = append(directions, rookDirections[:]...)

	// Normal directional moves
	for _, offset := range directions {
		pos := AddPositions(position, offset)
		if isInBounds(pos) &&
			isValidSquare(piece, state.Board[pos.X][pos.Y]) {

			possibleMove := Move{
				position,
				pos,
				None,
			}

			// Check that this move does not lead to the king being put into check
			if !IsSquareAttacked(state.Board.DoMove(possibleMove), pos, state.ActiveColor) {
				moves = append(moves, possibleMove)
			}
		}
	}

	// Castling
	if state.ActiveColor == Black && !IsSquareAttacked(state.Board, position, Black) {

		kingSidePositionFinish := Position{0, 6}
		kingSidePositionSkip := Position{0, 5}

		// Check if king side castling is legal (empty squares, no castling through check, etc.)
		if state.BlackCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, Black) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, Black) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
			})
		}

		queenSidePositionRookSkip := Position{0, 1}
		queenSidePositionFinish := Position{0, 2}
		queenSidePositionSkip := Position{0, 3}

		// Check if queen side castling is legal (empty squares, no castling through check, etc.)
		if state.BlackCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, Black) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, Black) {

			moves = append(moves, Move{
				position,
				queenSidePositionFinish,
				QueenSideCastle,
			})
		}
	} else if state.ActiveColor == White && !IsSquareAttacked(state.Board, position, White) {

		kingSidePositionFinish := Position{7, 6}
		kingSidePositionSkip := Position{7, 5}

		// Check if king side castling is legal (empty squares, no castling through check, etc.)
		if state.WhiteCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, White) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, White) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
			})
		}

		queenSidePositionRookSkip := Position{7, 1}
		queenSidePositionFinish := Position{7, 2}
		queenSidePositionSkip := Position{7, 3}

		// Check if queen side castling is legal (empty squares, no castling through check, etc.)
		if state.WhiteCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, White) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, White) {

			moves = append(moves, Move{
				position,
				queenSidePositionFinish,
				QueenSideCastle,
			})
		}
	}

	return
}

func GenerateQueenMoves(state State, position, kingPosition Position) (moves []Move) {
	moves = append(moves, generateDirectionalMoves(state, position, kingPosition, rookDirections)...)
	moves = append(moves, generateDirectionalMoves(state, position, kingPosition, bishopDirections)...)

	return
}

func GenerateRookMoves(state State, position, kingPosition Position) []Move {
	return generateDirectionalMoves(state, position, kingPosition, rookDirections)
}

func GenerateBishopMoves(state State, position, kingPosition Position) []Move {
	return generateDirectionalMoves(state, position, kingPosition, bishopDirections)
}

func GenerateKnightMoves(state State, position, kingPosition Position) (moves []Move) {
	checkIllegalMove := getIllegalMoveChecker(state, kingPosition)

	for _, knightPosition := range getKnightVision(position) {
		if isValidSquare(state.Board.GetSquare(position), state.Board.GetSquare(knightPosition)) {
			move := Move{
				position,
				knightPosition,
				None,
			}
			if !checkIllegalMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return
}

func GeneratePawnMoves(state State, position, kingPosition Position) (moves []Move) {

	isEnemyPiece := func(piecePosition Position) bool {
		piece := state.Board.GetSquare(piecePosition)
		return (state.ActiveColor == Black && piece < 0) ||
			(state.ActiveColor == White && piece > 0)
	}

	checkIllegalMove := getIllegalMoveChecker(state, kingPosition)

	isPromotion := false
	var pawnDirection int8
	var doublePushAvailable bool
	if state.ActiveColor == Black {
		pawnDirection = 1

		// Check if pawn is on the last rank
		if position.X == 6 {
			isPromotion = true
		}
		doublePushAvailable = position.X == 1
	} else {
		pawnDirection = -1

		// Check if pawn is on the last rank
		if position.X == 1 {
			isPromotion = true
		}
		doublePushAvailable = position.X == 6
	}

	// Normal moves
	pushPosition := AddPositions(position, Position{pawnDirection, 0})
	if isInBounds(pushPosition) && state.Board.GetSquare(pushPosition) == EmptySquare {
		move := Move{
			position,
			pushPosition,
			None,
		}

		if !checkIllegalMove(move) {
			if isPromotion {
				moves = append(moves, getMovesForPromotion(move)...)
			} else {
				moves = append(moves, move)
			}
		}

		// Double push
		if doublePushAvailable {
			doublePushPosition := AddPositions(position, Position{pawnDirection * 2, 0})
			if isInBounds(doublePushPosition) && state.Board.GetSquare(doublePushPosition) == EmptySquare {
				move = Move{
					position,
					doublePushPosition,
					None,
				}

				if !checkIllegalMove(move) {
					moves = append(moves, move)
				}
			}
		}
	}

	capturePositions := [2]Position{
		AddPositions(position, Position{pawnDirection, 1}),
		AddPositions(position, Position{pawnDirection, -1}),
	}

	for _, capturePosition := range capturePositions {
		if isInBounds(capturePosition) && isEnemyPiece(capturePosition) {
			move := Move{
				position,
				capturePosition,
				None,
			}

			if !checkIllegalMove(move) {
				if isPromotion {
					moves = append(moves, getMovesForPromotion(move)...)
				} else {
					moves = append(moves, move)
				}
			}
		} else if state.EnPassantPosition.Ok && capturePosition == state.EnPassantPosition.Position {
			move := Move{
				position,
				capturePosition,
				EnPassant,
			}

			if !checkIllegalMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return
}

// Checks if a square is attacked by the opposite color to defenderSide
func IsSquareAttacked(board Board, position Position, defenderSide Color) bool {

	// Check for attacks by pawn
	var pawnDirection int8
	var pawnType Piece
	if defenderSide == Black {
		pawnType = WhitePawn
		pawnDirection = 1
	} else if defenderSide == White {
		pawnType = BlackPawn
		pawnDirection = -1
	}

	pawnPositions := [2]Position{
		AddPositions(position, Position{pawnDirection, 1}),
		AddPositions(position, Position{pawnDirection, -1}),
	}

	for _, pawnPosition := range pawnPositions {
		if isInBounds(pawnPosition) &&
			board.GetSquare(pawnPosition) == pawnType {

			return true
		}
	}

	// Check for attacks by knight
	knightMoves := getKnightVision(position)
	for _, knightVision := range knightMoves {
		if (defenderSide == Black && board.GetSquare(knightVision) == WhiteKnight) ||
			(defenderSide == White && board.GetSquare(knightVision) == BlackKnight) {

			return true
		}
	}

	// Check for bishop type attacks
	bishopVisions := getDirectionalVision(board, position, bishopDirections)
	for _, bishopVision := range bishopVisions {
		// Don't consider the square itself
		if !bishopVision.Ok {
			continue
		}
		seenPiece := board.GetSquare(bishopVision.Position)
		if (defenderSide == Black && (seenPiece == WhiteBishop || seenPiece == WhiteQueen)) ||
			(defenderSide == White && (seenPiece == BlackBishop || seenPiece == BlackQueen)) {

			return true
		}
	}

	// Check for rook type attacks
	rookVisions := getDirectionalVision(board, position, rookDirections)
	for _, rookVision := range rookVisions {
		// Don't consider the square itself
		if !rookVision.Ok {
			continue
		}
		seenPiece := board.GetSquare(rookVision.Position)
		if (defenderSide == Black && (seenPiece == WhiteRook || seenPiece == WhiteQueen)) ||
			(defenderSide == White && (seenPiece == BlackRook || seenPiece == BlackQueen)) {

			return true
		}
	}

	// Check for attacks by king
	for _, offset := range rookDirections {
		singleRookMove := AddPositions(position, offset)
		if !isInBounds(singleRookMove) {
			continue
		}
		square := board.GetSquare(singleRookMove)
		if (defenderSide == Black && square == WhiteKing) ||
			(defenderSide == White && square == BlackKing) {

			return true
		}
	}
	for _, offset := range bishopDirections {
		singleBishopMove := AddPositions(position, offset)
		if !isInBounds(singleBishopMove) {
			continue
		}
		square := board.GetSquare(singleBishopMove)
		if (defenderSide == Black && square == WhiteKing) ||
			(defenderSide == White && square == BlackKing) {

			return true
		}
	}

	return false
}

// generateDirectionalMoves gives the available legal moves given a set of direction.
func generateDirectionalMoves(state State, position, kingPosition Position, directions [4]Position) (moves []Move) {
	checkIllegalMove := getIllegalMoveChecker(state, kingPosition)

	for i := 0; i < 4; i++ {
		offset := int8(1)
		for {
			nextPosition := AddPositions(position, MultiplyScalar(directions[i], offset))
			if !isInBounds(nextPosition) {
				break
			}

			possibleMove := Move{
				position,
				nextPosition,
				None,
			}

			// If the square is empty or an unfriendly piece and executing the move does not result in a check, add it to the list
			nextSquare := state.Board.GetSquare(nextPosition)
			isFriendlyPiece := (state.ActiveColor == Black && nextSquare > 0) || (state.ActiveColor == White && nextSquare < 0)
			if !isFriendlyPiece && !checkIllegalMove(possibleMove) {
				moves = append(moves, possibleMove)
			}

			// Break on any non-empty squares
			if nextSquare != EmptySquare {
				break
			}
			offset++
		}
	}

	return
}

// getDirectionalVision gives the extents of the vision/movement of a piece.
func getDirectionalVision(board Board, position Position, directions [4]Position) [4]PositionOpt {
	optPositions := [4]PositionOpt{}

	for i := 0; i < 4; i++ {
		endFound := false
		offset := int8(1)
		for !endFound {
			nextPosition := AddPositions(position, MultiplyScalar(directions[i], offset))
			if !isInBounds(nextPosition) {
				endFound = true
				optPositions[i] = PositionOpt{Ok: false}
			} else if nextSquare := board.GetSquare(nextPosition); nextSquare != EmptySquare {
				endFound = true
				optPositions[i] = PositionOpt{
					Position: nextPosition,
					Ok:       true,
				}
			}
			offset++
		}
	}

	return optPositions
}

// getKnightVision returns all possible knight jumps given a position.
func getKnightVision(position Position) []Position {
	validPositions := []Position{}
	for _, offset := range knightOffsets {
		possiblePosition := AddPositions(position, offset)
		if isInBounds(possiblePosition) {
			validPositions = append(validPositions, possiblePosition)
		}
	}

	return validPositions
}

// findKing returns the position of the king with the specified color.
func findKing(board Board, color Color) (kingPosition Position, err error) {
	var square Piece

	var i int8
	var j int8
	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			kingPosition = Position{i, j}
			square = board.GetSquare(kingPosition)

			if (color == Black && square == BlackKing) ||
				(color == White && square == WhiteKing) {

				return kingPosition, nil
			}
		}
	}

	return Position{}, errors.New("missing king")
}

// isInCheck returns whether the given color's king is in check.
func isInCheck(board Board, color Color) (bool, error) {
	kingPosition, err := findKing(board, color)
	if err != nil {
		return false, err
	}

	return IsSquareAttacked(board, kingPosition, color), nil
}

// isInBounds checks if the position falls within the board.
func isInBounds(position Position) bool {
	return position.X >= 0 && position.X <= 7 && position.Y >= 0 && position.Y <= 7
}

// isValidSquare checks if the square is empty or an enemy piece.
func isValidSquare(piece1, piece2 Piece) bool {
	return piece1*piece2 <= 0
}

// getIllegalMoveChecker returns a method that checks whether a given move is legal with the current state and king square.
func getIllegalMoveChecker(state State, kingPosition Position) func(move Move) bool {
	return func(move Move) bool {
		return IsSquareAttacked(state.Board.DoMove(move), kingPosition, state.ActiveColor)
	}
}

// getMovesForPromotion takes the positional information from the parameter and returns an array with each type of promotion.
func getMovesForPromotion(move Move) []Move {
	moves := make([]Move, 4)

	moves[0] = Move{
		move.Start,
		move.End,
		PromoteToQueen,
	}
	moves[1] = Move{
		move.Start,
		move.End,
		PromoteToRook,
	}
	moves[2] = Move{
		move.Start,
		move.End,
		PromoteToBishop,
	}
	moves[3] = Move{
		move.Start,
		move.End,
		PromoteToKnight,
	}

	return moves
}
