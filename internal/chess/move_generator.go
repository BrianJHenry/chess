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
func (state *State) GenerateAllMoves() (moves []Move, err error) {
	kingPosition, err := state.Board.findKing(state.ActiveColor)
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
				moves = append(moves, state.GeneratePawnMoves(position, kingPosition)...)
			case WhiteBishop, BlackBishop:
				moves = append(moves, state.GenerateBishopMoves(position, kingPosition)...)
			case WhiteKnight, BlackKnight:
				moves = append(moves, state.GenerateKnightMoves(position, kingPosition)...)
			case WhiteRook, BlackRook:
				moves = append(moves, state.GenerateRookMoves(position, kingPosition)...)
			case WhiteQueen, BlackQueen:
				moves = append(moves, state.GenerateQueenMoves(position, kingPosition)...)
			case WhiteKing, BlackKing:
				moves = append(moves, state.GenerateKingMoves(position)...)
			default:
				// Pass on empty squares
			}
		}
	}

	return
}

func (state *State) GenerateKingMoves(position Position) (moves []Move) {
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
				state.Board[pos.X][pos.Y],
			}

			// Check that this move does not lead to the king being put into check
			state.Board.DoMove(possibleMove)
			if !state.Board.IsSquareAttacked(pos, state.ActiveColor) {
				moves = append(moves, possibleMove)
			}
			state.Board.UndoMove(possibleMove)
		}
	}

	// Castling
	if state.ActiveColor == Black && !state.Board.IsSquareAttacked(position, Black) {

		kingSidePositionFinish := Position{0, 6}
		kingSidePositionSkip := Position{0, 5}

		// Check if king side castling is legal (empty squares, no castling through check, etc.)
		if state.CastlingRights.BlackCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!state.Board.IsSquareAttacked(kingSidePositionSkip, Black) &&
			!state.Board.IsSquareAttacked(kingSidePositionFinish, Black) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
				EmptySquare,
			})
		}

		queenSidePositionRookSkip := Position{0, 1}
		queenSidePositionFinish := Position{0, 2}
		queenSidePositionSkip := Position{0, 3}

		// Check if queen side castling is legal (empty squares, no castling through check, etc.)
		if state.CastlingRights.BlackCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!state.Board.IsSquareAttacked(queenSidePositionSkip, Black) &&
			!state.Board.IsSquareAttacked(queenSidePositionFinish, Black) {

			moves = append(moves, Move{
				position,
				queenSidePositionFinish,
				QueenSideCastle,
				EmptySquare,
			})
		}
	} else if state.ActiveColor == White && !state.Board.IsSquareAttacked(position, White) {

		kingSidePositionFinish := Position{7, 6}
		kingSidePositionSkip := Position{7, 5}

		// Check if king side castling is legal (empty squares, no castling through check, etc.)
		if state.CastlingRights.WhiteCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!state.Board.IsSquareAttacked(kingSidePositionSkip, White) &&
			!state.Board.IsSquareAttacked(kingSidePositionFinish, White) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
				EmptySquare,
			})
		}

		queenSidePositionRookSkip := Position{7, 1}
		queenSidePositionFinish := Position{7, 2}
		queenSidePositionSkip := Position{7, 3}

		// Check if queen side castling is legal (empty squares, no castling through check, etc.)
		if state.CastlingRights.WhiteCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!state.Board.IsSquareAttacked(queenSidePositionSkip, White) &&
			!state.Board.IsSquareAttacked(queenSidePositionFinish, White) {

			moves = append(moves, Move{
				position,
				queenSidePositionFinish,
				QueenSideCastle,
				EmptySquare,
			})
		}
	}

	return
}

func (state *State) GenerateQueenMoves(position, kingPosition Position) (moves []Move) {
	moves = append(moves, state.generateDirectionalMoves(position, kingPosition, rookDirections)...)
	moves = append(moves, state.generateDirectionalMoves(position, kingPosition, bishopDirections)...)

	return
}

func (state *State) GenerateRookMoves(position, kingPosition Position) []Move {
	return state.generateDirectionalMoves(position, kingPosition, rookDirections)
}

func (state *State) GenerateBishopMoves(position, kingPosition Position) []Move {
	return state.generateDirectionalMoves(position, kingPosition, bishopDirections)
}

func (state *State) GenerateKnightMoves(position, kingPosition Position) (moves []Move) {
	checkIllegalMove := state.getIllegalMoveChecker(kingPosition)

	for _, knightPosition := range getKnightVision(position) {
		if isValidSquare(state.Board.GetSquare(position), state.Board.GetSquare(knightPosition)) {
			move := Move{
				position,
				knightPosition,
				None,
				state.Board.GetSquare((knightPosition)),
			}
			if !checkIllegalMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return
}

func (state *State) GeneratePawnMoves(position, kingPosition Position) (moves []Move) {

	isEnemyPiece := func(piecePosition Position) bool {
		piece := state.Board.GetSquare(piecePosition)
		return (state.ActiveColor == Black && piece < 0) ||
			(state.ActiveColor == White && piece > 0)
	}

	checkIllegalMove := state.getIllegalMoveChecker(kingPosition)

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
			EmptySquare,
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
					EmptySquare,
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
				state.Board.GetSquare(capturePosition),
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
				getEnemyPawnColor(state.ActiveColor),
			}

			if !checkIllegalMove(move) {
				moves = append(moves, move)
			}
		}
	}

	return
}

// Checks if a square is attacked by the opposite color to defenderSide
func (board *Board) IsSquareAttacked(position Position, defenderSide Color) bool {

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
	bishopVisions := board.getDirectionalVision(position, bishopDirections)
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
	rookVisions := board.getDirectionalVision(position, rookDirections)
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
func (state *State) generateDirectionalMoves(position, kingPosition Position, directions [4]Position) (moves []Move) {
	checkIllegalMove := state.getIllegalMoveChecker(kingPosition)

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
				state.Board.GetSquare(nextPosition),
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
func (board *Board) getDirectionalVision(position Position, directions [4]Position) [4]PositionOpt {
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
func (board *Board) findKing(color Color) (kingPosition Position, err error) {
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
func (board *Board) isInCheck(color Color) (bool, error) {
	kingPosition, err := board.findKing(color)
	if err != nil {
		return false, err
	}

	return board.IsSquareAttacked(kingPosition, color), nil
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
func (state *State) getIllegalMoveChecker(kingPosition Position) func(move Move) bool {
	return func(move Move) bool {
		state.Board.DoMove(move)
		isAttacked := state.Board.IsSquareAttacked(kingPosition, state.ActiveColor)
		state.Board.UndoMove(move)
		return isAttacked
	}
}

// getMovesForPromotion takes the positional information from the parameter and returns an array with each type of promotion.
func getMovesForPromotion(move Move) []Move {
	moves := make([]Move, 4)

	moves[0] = Move{
		move.Start,
		move.End,
		PromoteToQueen,
		move.Captured,
	}
	moves[1] = Move{
		move.Start,
		move.End,
		PromoteToRook,
		move.Captured,
	}
	moves[2] = Move{
		move.Start,
		move.End,
		PromoteToBishop,
		move.Captured,
	}
	moves[3] = Move{
		move.Start,
		move.End,
		PromoteToKnight,
		move.Captured,
	}

	return moves
}
