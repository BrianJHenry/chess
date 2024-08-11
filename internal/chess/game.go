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

	possibleMoves, err := initialState.GenerateAllMoves()
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
func (game *Game) DoMove(move Move) error {
	game.State.DoMove(move)
	game.Moves = append(game.Moves, move)

	possibleMoves, err := game.State.GenerateAllMoves()
	if err != nil {
		return err
	}

	checkmate := false
	stalemate := false
	if len(possibleMoves) == 0 {
		isInCheck, err := game.State.Board.isInCheck(game.State.ActiveColor)
		if err != nil {
			return err
		}

		if isInCheck {
			checkmate = true
		} else {
			stalemate = false
		}
	}

	game.Checkmate = checkmate
	game.Stalemate = stalemate
	return nil
}
