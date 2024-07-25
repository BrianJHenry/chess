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

func GenerateAllMoves(state State) (moves []Move, err error) {
	moves = []Move{}

	var position Position
	var square Piece
	var i int8
	var j int8

	// TODO: improve this!
	kingPosition, err := state.Board.FindKing(state.Turn)
	if err != nil {
		return moves, err
	}

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			position = Position{i, j}
			square = state.Board.GetSquare(position)

			if (state.Turn == BlackTurn && square <= 0) || (state.Turn == WhiteTurn && square >= 0) {
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
	moves = []Move{}

	piece := state.Board[position.X][position.Y]

	// Normal moves
	for _, offset := range bishopDirections {
		pos := position.AddOffset(offset)
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

	for _, offset := range rookDirections {
		pos := position.AddOffset(offset)
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

		kingSidePositionFinish := Position{0, 6}
		kingSidePositionSkip := Position{0, 5}

		if state.BlackCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, BlackTurn) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, BlackTurn) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
			})
		}

		queenSidePositionRookSkip := Position{0, 1}
		queenSidePositionFinish := Position{0, 2}
		queenSidePositionSkip := Position{0, 3}

		if state.BlackCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, BlackTurn) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, BlackTurn) {

			moves = append(moves, Move{
				position,
				queenSidePositionFinish,
				QueenSideCastle,
			})
		}
	} else {

		kingSidePositionFinish := Position{7, 6}
		kingSidePositionSkip := Position{7, 5}

		if state.WhiteCanCastleKingSide &&
			state.Board.GetSquare(kingSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(kingSidePositionFinish) == EmptySquare &&
			!IsSquareAttacked(state.Board, kingSidePositionSkip, WhiteTurn) &&
			!IsSquareAttacked(state.Board, kingSidePositionFinish, WhiteTurn) {

			moves = append(moves, Move{
				position,
				kingSidePositionFinish,
				KingSideCastle,
			})
		}

		queenSidePositionRookSkip := Position{7, 1}
		queenSidePositionFinish := Position{7, 2}
		queenSidePositionSkip := Position{7, 3}

		if state.WhiteCanCastleQueenSide &&
			state.Board.GetSquare(queenSidePositionSkip) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionFinish) == EmptySquare &&
			state.Board.GetSquare(queenSidePositionRookSkip) == EmptySquare &&
			!IsSquareAttacked(state.Board, queenSidePositionSkip, WhiteTurn) &&
			!IsSquareAttacked(state.Board, queenSidePositionFinish, WhiteTurn) {

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
	moves = []Move{}

	moves = append(moves, generateDirectionalMoves(state, position, kingPosition, rookDirections)...)
	moves = append(moves, generateDirectionalMoves(state, position, kingPosition, bishopDirections)...)

	return
}

func GenerateRookMoves(state State, position, kingPosition Position) (moves []Move) {
	moves = generateDirectionalMoves(state, position, kingPosition, rookDirections)
	return
}

func GenerateBishopMoves(state State, position, kingPosition Position) (moves []Move) {
	moves = generateDirectionalMoves(state, position, kingPosition, bishopDirections)
	return
}

func GenerateKnightMoves(state State, position, kingPosition Position) (moves []Move) {
	moves = []Move{}

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
		return (state.Turn == BlackTurn && state.Board.GetSquare(piecePosition) < 0) ||
			(state.Turn == WhiteTurn && state.Board.GetSquare(piecePosition) > 0)
	}

	checkIllegalMove := getIllegalMoveChecker(state, kingPosition)

	isPromotion := false
	var pawnDirection int8
	var doublePushAvailable bool
	if state.Turn == BlackTurn {
		pawnDirection = 1
		if position.X == 6 {
			isPromotion = true
		}
		doublePushAvailable = position.X == 1
	} else {
		pawnDirection = -1
		if position.X == 1 {
			isPromotion = true
		}
		doublePushAvailable = position.X == 6
	}

	// Normal moves
	pushPosition := position.AddOffset(Position{pawnDirection, 0})
	var move Move
	if isInBounds(pushPosition) && state.Board.GetSquare(pushPosition) == EmptySquare {
		move = Move{
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
			doublePushPosition := position.AddOffset(Position{pawnDirection * 2, 0})
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
		position.AddOffset(Position{pawnDirection, 1}),
		position.AddOffset(Position{pawnDirection, -1}),
	}

	for _, capturePosition := range capturePositions {
		if isInBounds(capturePosition) && isEnemyPiece(capturePosition) {
			move = Move{
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
		} else if state.EnPassantPosition.Valid && capturePosition == state.EnPassantPosition.Position {
			move = Move{
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
func IsSquareAttacked(board Board, position Position, defenderSide Turn) bool {

	// Check for attacks by pawn
	var pawnDirection int8
	var pawnType Piece
	if defenderSide == BlackTurn {
		pawnType = WhitePawn
		pawnDirection = 1
	} else if defenderSide == WhiteTurn {
		pawnType = BlackPawn
		pawnDirection = -1
	}

	pawnPositions := [2]Position{
		position.AddOffset(Position{pawnDirection, 1}),
		position.AddOffset(Position{pawnDirection, -1}),
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
		if (defenderSide == BlackTurn && board.GetSquare(knightVision) == WhiteKnight) ||
			(defenderSide == WhiteTurn && board.GetSquare(knightVision) == BlackKnight) {

			return true
		}
	}

	// Check for bishop type attacks
	bishopVisions := getDirectionalVision(board, position, bishopDirections)
	for _, bishopVision := range bishopVisions {
		// Don't consider the square itself
		if position == bishopVision {
			continue
		}
		seenPiece := board.GetSquare(bishopVision)
		if (defenderSide == BlackTurn && (seenPiece == WhiteBishop || seenPiece == WhiteQueen)) ||
			(defenderSide == WhiteTurn && (seenPiece == BlackBishop || seenPiece == BlackQueen)) {

			return true
		}
	}

	// Check for rook type attacks
	rookVisions := getDirectionalVision(board, position, rookDirections)
	for _, rookVision := range rookVisions {
		// Don't consider the square itself
		if position == rookVision {
			continue
		}
		seenPiece := board.GetSquare(rookVision)
		if (defenderSide == BlackTurn && (seenPiece == WhiteRook || seenPiece == WhiteQueen)) ||
			(defenderSide == WhiteTurn && (seenPiece == BlackRook || seenPiece == BlackQueen)) {

			return true
		}
	}

	// Check for attacks by king
	var square Piece
	for _, offset := range rookDirections {
		singleRookMove := position.AddOffset(offset)
		if !isInBounds(singleRookMove) {
			continue
		}
		square = board.GetSquare(singleRookMove)
		if (defenderSide == BlackTurn && square == WhiteKing) ||
			(defenderSide == WhiteTurn && square == BlackKing) {

			return true
		}
	}

	for _, offset := range bishopDirections {
		singleBishopMove := position.AddOffset(offset)
		if !isInBounds(singleBishopMove) {
			continue
		}
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

func generateDirectionalMoves(state State, position, kingPosition Position, directions [4]Position) (moves []Move) {
	moves = []Move{}

	checkIllegalMove := getIllegalMoveChecker(state, kingPosition)

	var offset int8
	var possibleMove Move
	for i := 0; i < 4; i++ {
		offset = 1
		for {
			nextPosition := position.AddOffset(directions[i].MultiplyScalar(offset))

			if !isInBounds(nextPosition) {
				break
			}

			nextSquare := state.Board.GetSquare(nextPosition)
			possibleMove = Move{
				position,
				nextPosition,
				None,
			}

			// If the square is empty or an unfriendly piece and executing the move does not result in a check, add it to the list
			isFriendlyPiece := (state.Turn == BlackTurn && nextSquare > 0) || (state.Turn == WhiteTurn && nextSquare < 0)
			if !isFriendlyPiece && !checkIllegalMove(possibleMove) {
				moves = append(moves, possibleMove)
			}

			// Break on any non-empty squares
			if nextSquare != 0 {
				break
			}
			offset++
		}
	}

	return
}

// The last square in a set of directions
// TODO: Possibly return NullablePosition (case where a piece is up against the side of the board)
func getDirectionalVision(board Board, position Position, directions [4]Position) [4]Position {
	positions := [4]Position{}

	var endFound bool
	var offset int8
	for i := 0; i < 4; i++ {
		endFound = false
		offset = 1
		for !endFound {
			nextPosition := position.AddOffset(directions[i].MultiplyScalar(offset))
			if !isInBounds(nextPosition) {
				endFound = true
				positions[i] = position.AddOffset(directions[i].MultiplyScalar(offset - 1))
			} else if nextSquare := board.GetSquare(nextPosition); nextSquare != EmptySquare {
				endFound = true
				positions[i] = nextPosition
			}
			offset++
		}
	}

	return positions
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

func getIllegalMoveChecker(state State, kingPosition Position) func(move Move) bool {
	return func(move Move) bool {
		return IsSquareAttacked(state.Board.ExecuteMove(move), kingPosition, state.Turn)
	}
}

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
