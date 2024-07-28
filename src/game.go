package chess

type Game struct {
	State         State
	Moves         []Move
	Checkmate     bool
	Stalemate     bool
	possibleMoves []Move
}

// DoMove takes in a Game object and a Move and executes the move, returning the updated Game object.
func (game Game) DoMove(move Move) (Game, error) {
	state := game.State.DoMove(move)
	moves := append(game.Moves, move)

	possibleMoves, err := GenerateAllMoves(state)
	if err != nil {
		return Game{}, err
	}

	checkmate := false
	stalemate := false
	if len(possibleMoves) == 0 {
		isInCheck, err := IsInCheck(state.Board, state.ActiveColor)
		if err != nil {
			return Game{}, err
		}

		if isInCheck {
			checkmate = true
		} else {
			stalemate = false
		}
	}

	return Game{
		state,
		moves,
		checkmate,
		stalemate,
		possibleMoves,
	}, nil
}
