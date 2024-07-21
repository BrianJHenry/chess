package chess

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

// TDOO
func GenerateAllMoves(state State) (moves []Move) {
	return
}

func GenerateKingMoves(state State, position Position) (moves []Move) {
	moves = []Move{}

	piece := state.Board[position.X][position.Y]

	// Normal moves
	for _, pos := range bishopDirections {
		if isInBounds(pos) &&
			isValidSquare(piece, state.Board[pos.X][pos.Y]) &&
			!IsSquareAttacked(state.Board, pos, state.Turn) {
			moves = append(moves, Move{
				position,
				pos,
				None,
			})
		}
	}

	for _, pos := range rookDirections {
		if isInBounds(pos) &&
			isValidSquare(piece, state.Board[pos.X][pos.Y]) &&
			!IsSquareAttacked(state.Board, pos, state.Turn) {
			moves = append(moves, Move{
				position,
				pos,
				None,
			})
		}
	}

	// Castling
	if state.Turn == BlackTurn {

		kingSidePositionFinish := Position{X: 1, Y: 7}
		kingSidePositionSkip := Position{X: 2, Y: 7}

		if state.BlackCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, BlackTurn) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, BlackTurn) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				Castle,
			})
		}

		queenSidePositionRookSkip := Position{X: 6, Y: 7}
		queenSidePositionFinish := Position{X: 5, Y: 7}
		queenSidePositionSkip := Position{X: 4, Y: 7}

		if state.BlackCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, BlackTurn) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, BlackTurn) {

			moves = append(moves, Move{
				position,
				Position{X: 5, Y: 7},
				Castle,
			})
		}
	} else {

		kingSidePositionFinish := Position{X: 1, Y: 0}
		kingSidePositionSkip := Position{X: 2, Y: 0}

		if state.WhiteCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, WhiteTurn) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, WhiteTurn) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				Castle,
			})
		}

		queenSidePositionRookSkip := Position{X: 6, Y: 0}
		queenSidePositionFinish := Position{X: 5, Y: 0}
		queenSidePositionSkip := Position{X: 4, Y: 0}

		if state.WhiteCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, WhiteTurn) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, WhiteTurn) {

			moves = append(moves, Move{
				position,
				Position{X: 5, Y: 0},
				Castle,
			})
		}
	}

	return
}

// TODO
func GenerateQueenMoves(state State, position, kingPosition Position) (moves []Move) {
	return
}

// TODO
func GenerateRookMoves(state State, position, kingPosition Position) (moves []Move) {
	return
}

// TODO
func GenerateBishopMoves(state State, position, kingPosition Position) (moves []Move) {
	return
}

// TODO
func GenerateKnightMoves(state State, position, kingPosition Position) (moves []Move) {
	return
}

// TODO
func GeneratePawnMoves(state State, position, kingPosition Position) (moves []Move) {
	return
}

// Checks if a square is attacked by the opposite color to defenderSide
func IsSquareAttacked(board Board, position Position, defenderSide Turn) bool {

	// Check for attacks by pawn
	var pawnDirection int8
	var pawnType Piece
	if defenderSide == BlackTurn {
		pawnType = WhitePawn
		pawnDirection = -1
	} else if defenderSide == WhiteTurn {
		pawnType = BlackPawn
		pawnDirection = 1
	}

	pawnPositions := [2]Position{
		position.AddOffset(Position{X: 1, Y: pawnDirection}),
		position.AddOffset(Position{X: -1, Y: pawnDirection}),
	}

	for _, pawnPosition := range pawnPositions {
		if isInBounds(pawnPosition) &&
			board.GetSquare(pawnPosition) == pawnType {

			return true
		}
	}

	// Check for attacks by knight
	for _, knightVision := range getKnightVision(position) {
		if defenderSide == BlackTurn && board.GetSquare(knightVision) == WhiteKnight {
			return true
		} else if defenderSide == WhiteTurn && board.GetSquare(knightVision) == BlackKnight {
			return true
		}
	}

	// Check for bishop type attacks
	for _, bishopVision := range getDirectionalVision(board, position, bishopDirections) {
		if (defenderSide == BlackTurn && (bishopVision == WhiteBishop || bishopVision == WhiteQueen)) ||
			(defenderSide == WhiteTurn && (bishopVision == BlackBishop || bishopVision == BlackQueen)) {

			return true
		}
	}

	// Check for rook type attacks
	for _, rookVision := range getDirectionalVision(board, position, rookDirections) {
		if (defenderSide == BlackTurn && (rookVision == WhiteRook || rookVision == WhiteQueen)) ||
			(defenderSide == WhiteTurn && (rookVision == BlackRook || rookVision == BlackQueen)) {

			return true
		}
	}

	// Check for attacks by king
	var square Piece
	for _, singleRookMove := range rookDirections {
		square = board.GetSquare(singleRookMove)
		if (defenderSide == BlackTurn && square == WhiteKing) ||
			(defenderSide == WhiteTurn && square == BlackKing) {

			return true
		}
	}

	for _, singleBishopMove := range bishopDirections {
		square = board.GetSquare(singleBishopMove)
		if (defenderSide == BlackTurn && square == WhiteKing) ||
			(defenderSide == WhiteTurn && square == BlackKing) {

			return true
		}
	}

	return false
}

func isInBounds(position Position) bool {
	return position.X >= 0 && position.X <= 7 && position.Y >= 0 && position.Y <= 7
}

func isValidSquare(piece1, piece2 Piece) bool {
	return piece1*piece2 <= 0
}

// The last square in a set of directions
func getDirectionalVision(board Board, position Position, directions [4]Position) [4]Piece {
	pieces := [4]Piece{}

	var endFound bool
	var offset int8
	for i := 0; i < 4; i++ {
		endFound = false
		offset = 1
		for !endFound {
			nextPosition := position.AddOffset(directions[i].MultiplyScalar(offset))
			if !isInBounds(nextPosition) {
				endFound = true
				pieces[i] = EmptySquare
			} else if nextSquare := board.GetSquare(nextPosition); nextSquare != EmptySquare {
				endFound = true
				pieces[i] = nextSquare
			}
		}
	}

	return pieces
}

func getKnightVision(position Position) []Position {
	validPositions := []Position{}
	for _, offset := range knightOffsets {
		possiblePosition := position.AddOffset(offset)
		if isInBounds(possiblePosition) {
			validPositions = append(validPositions, possiblePosition)
		}
	}

	return validPositions
}
