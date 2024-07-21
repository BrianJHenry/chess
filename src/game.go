package chess

type Game struct {
	State         State
	Moves         []Move
	Checkmate     bool
	Stalemate     bool
	possibleMoves []Move
}

func (game Game) ExecuteMove(move Move) (Game, error) {
	state := game.State.ExecuteMove(move)
	moves := append(game.Moves, move)

	possibleMoves, err := GenerateAllMoves(state)
	if err != nil {
		return Game{}, err
	}

	checkmate := false
	stalemate := false
	if len(possibleMoves) == 0 {
		isInCheck, err := state.Board.IsInCheck(state.Turn)
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
