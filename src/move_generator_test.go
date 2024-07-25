package chess

import (
	"errors"
	"fmt"
	"testing"
)

func TestMoveGenerationAndExecution(t *testing.T) {
	data, err := LoadTestData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	errorCollection := []error{}
	for _, test := range data {
		moves, err := GenerateAllMoves(test.Initial)
		if err != nil {
			errorCollection = append(errorCollection, err)
			continue
		}

		nextStates := make([]State, len(moves))
		for i, move := range moves {
			nextStates[i] = test.Initial.ExecuteMove(move)
		}

		// TODO: come back and refine data structure/comparison; N^2 is ugly
		for _, result := range test.Results {
			matched := false
			for _, nextState := range nextStates {
				if result.Result.Equals(nextState) {
					matched = true
					break
				}
			}

			if !matched {
				errorCollection = append(errorCollection, fmt.Errorf("move not found in state %s: %s\nInitial:\n%s\nResult\n%s", test.Description, result.Description, test.Initial.Board.GetPrintableBoard(), result.Result.Board.GetPrintableBoard()))
			}
		}
	}

	err = errors.Join(errorCollection...)
	if err != nil {
		t.Fatalf("Errors: %d\n%s", len(errorCollection), err.Error())
	}
}
