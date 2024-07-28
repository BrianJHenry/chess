package engine

import "github.com/BrianJHenry/chess/internal/chess"

type Engine interface {
	ChooseMove(state chess.State, options []chess.Move) chess.Move
}

func ResolveEngine(engineVersion string) (Engine, error) {
	switch engineVersion {
	case "random":
		return RandomEngine{}, nil
	default:
		// TODO: figure out default engine
		return RandomEngine{}, nil
	}
}
