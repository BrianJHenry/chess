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
			position = Position{X: i, Y: j}
			square = state.Board.GetSquare(position)

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
				KingSideCastle,
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
				QueenSideCastle,
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
				KingSideCastle,
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
			if checkIllegalMove(move) {
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

	var pawnDirection int8
	var promotionFlag MoveFlag
	var doublePushAvailable bool
	if state.Turn == BlackTurn {
		pawnDirection = 1
		if position.X == 6 {
			promotionFlag = PromoteToQueen
		} else {
			promotionFlag = None
		}
		doublePushAvailable = position.X == 1
	} else {
		pawnDirection = -1
		if position.X == 1 {
			promotionFlag = PromoteToQueen
		} else {
			promotionFlag = None
		}
		doublePushAvailable = position.X == 6
	}

	// Normal moves
	pushPosition := position.AddOffset(Position{X: 0, Y: pawnDirection})
	var move Move
	if isInBounds(pushPosition) && state.Board.GetSquare(pushPosition) == EmptySquare {
		move = Move{
			position,
			pushPosition,
			promotionFlag,
		}

		if !checkIllegalMove(move) {
			moves = append(moves, move)
		}

		// Double push
		if doublePushAvailable {
			doublePushPosition := position.AddOffset(Position{X: 0, Y: pawnDirection * 2})
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
		position.AddOffset(Position{X: 1, Y: pawnDirection}),
		position.AddOffset(Position{X: -1, Y: pawnDirection}),
	}

	for _, capturePosition := range capturePositions {
		if isInBounds(capturePosition) && isEnemyPiece(capturePosition) {
			move = Move{
				position,
				capturePosition,
				promotionFlag,
			}

			if !checkIllegalMove(move) {
				moves = append(moves, move)
			}
		}
	}

	// En Passant
	previousMove := state.Previous.DecodeMove()
	previousMovePiece := state.Board.GetSquare(previousMove.End)
	previousDoublePawnPush := (previousMovePiece == BlackPawn && previousMove.Start.Y == 6 && previousMove.End.Y == 4) ||
		(previousMovePiece == WhitePawn && previousMove.Start.Y == 1 && previousMove.End.Y == 3)

	if previousDoublePawnPush && position.Y == previousMove.End.Y &&
		(position.X+1 == previousMove.End.X || position.X-1 == previousMove.End.X) {

		move = Move{
			position,
			Position{X: previousMove.End.X, Y: position.Y + pawnDirection},
			EnPassant,
		}

		if !checkIllegalMove(move) {
			moves = append(moves, move)
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
		}
	}

	return
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

func getIllegalMoveChecker(state State, kingPosition Position) func(move Move) bool {
	return func(move Move) bool {
		return IsSquareAttacked(state.Board.ExecuteMove(move), kingPosition, state.Turn)
	}
}
