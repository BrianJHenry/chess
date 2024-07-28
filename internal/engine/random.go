package engine

import (
	"math/rand"

	"github.com/BrianJHenry/chess/internal/chess"
)

type RandomEngine struct{}

func (e RandomEngine) ChooseMove(state chess.State, options []chess.Move) chess.Move {
	r := rand.New(rand.NewSource(0))

	return options[r.Int()%len(options)]
}
