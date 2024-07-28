package chess

type Game struct {
	State         State
	Moves         []Move
	Checkmate     bool
	Stalemate     bool
	PossibleMoves []Move
}

func InitialiseGame() Game {
	initialState := InitialiseState()

	possibleMoves, err := GenerateAllMoves(initialState)
	if err != nil {
		panic("call to GenerateAllMoves should never fail in the opening position.")
	}

	return Game{
		State:         initialState,
		Moves:         []Move{},
		Checkmate:     false,
		Stalemate:     false,
		PossibleMoves: possibleMoves,
	}
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
		isInCheck, err := isInCheck(state.Board, state.ActiveColor)
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
