package chess

import (
	"errors"
	"fmt"
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

	errorCollection := []error{}
	for _, test := range data {
		moves, err := GenerateAllMoves(test.Initial)
		if err != nil {
			errorCollection = append(errorCollection, err)
			continue
		}

		nextStates := make([]StateAndPreviousMove, len(moves))
		for i, move := range moves {
			nextState := test.Initial.ExecuteMove(move)
			moveString, err := move.ToAlgebraicNotation(test.Initial)
			if err != nil {
				errorCollection = append(errorCollection, err)
			}

			nextStates[i] = StateAndPreviousMove{
				State: nextState,
				Move:  moveString,
			}
		}

		// TODO: come back and refine data structure/comparison; N^2 is ugly
		for _, result := range test.Results {
			isBoardMatched := false
			for _, nextState := range nextStates {
				if result.Result == nextState.State {
					isBoardMatched = true

					// Check that move notation is correct
					if result.Move != nextState.Move {
						errorCollection = append(errorCollection, fmt.Errorf("move notation incorrect: expected=%s ; actual=%s", result.Move, nextState.Move))
					}
					break
				}
			}

			if !isBoardMatched {
				errorCollection = append(errorCollection, fmt.Errorf("state not found for move %s: %s\nInitial:\n%s\nResult\n%s", test.Description, result.Move, test.Initial.Board.GetPrintableBoard(), result.Result.Board.GetPrintableBoard()))
			}

		}
	}

	err = errors.Join(errorCollection...)
	if err != nil {
		t.Fatalf("Errors: %d\n%s", len(errorCollection), err.Error())
	}
}
