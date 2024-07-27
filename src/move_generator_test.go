package chess

import (
	"testing"
)

type StateAndPreviousMove struct {
	State State
	Move  string
}

func TestMoveGenerationAndExecution(t *testing.T) {
	data, err := LoadTestData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, test := range data {
		moves, err := GenerateAllMoves(test.Initial)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		// Execute and create move lookup
		moveLookup := make(map[string]Move, len(moves))
		for _, move := range moves {
			algebraic, err := move.ToAlgebraicNotation(test.Initial)
			if err != nil {
				t.Errorf(err.Error())
				continue
			}
			moveLookup[string(algebraic)] = move
		}

		// Check that all moves were properly found
		for _, result := range test.Results {
			move, ok := moveLookup[result.Move]
			if !ok {
				t.Errorf("move not generated %s", result.Move)
			}

			generatedState := test.Initial.ExecuteMove(move)
			if generatedState != result.Result {
				t.Errorf("incorrect resultant state for move %s\nexpected=\n%s\nactual=\n%s", result.Move, result.Result.Board.GetPrintableBoard(), generatedState.Board.GetPrintableBoard())
			}
		}
	}
}
