package chess

import (
	"testing"
)

type StateAndPreviousMove struct {
	State State
	Move  string
}

func TestMoveGenerationAndExecution(t *testing.T) {
	data, err := loadTestData()
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, test := range data {
		moves, err := test.Initial.GenerateAllMoves()
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

		if len(moveLookup) != len(test.Results) {
			t.Errorf("unexpected move count: expected=%d; actual=%d", len(test.Results), len(moveLookup))
		}

		// Check that all moves were properly found
		for _, result := range test.Results {
			move, ok := moveLookup[result.Move]
			if !ok {
				t.Errorf("move not generated %s", result.Move)
			}

			castlingRights := test.Initial.CastlingRights
			enPassantSquare := test.Initial.EnPassantPosition

			test.Initial.DoMove(move)
			if test.Initial != result.Result {
				t.Errorf("incorrect resultant state for move %s\nexpected=\n%s\nactual=\n%s", result.Move, BoardToDisplayString(result.Result.Board), BoardToDisplayString(test.Initial.Board))
			}
			test.Initial.UndoMove(move, castlingRights, enPassantSquare)
		}
	}
}

func TestMoveCountStartingPositionDepth1(t *testing.T) {
	testMoveCount(t, InitialiseState(), 20, 1)
}

func TestMoveCountStartingPositionDepth2(t *testing.T) {
	testMoveCount(t, InitialiseState(), 400, 2)
}

func TestMoveCountStartingPositionDepth3(t *testing.T) {
	testMoveCount(t, InitialiseState(), 8902, 3)
}

func TestMoveCountStartingPositionDepth4(t *testing.T) {
	testMoveCount(t, InitialiseState(), 197281, 4)
}

func TestMoveCountStartingPositionDepth5(t *testing.T) {
	testMoveCount(t, InitialiseState(), 4865609, 5)
}

func TestMoveCountStartingPositionDepth6(t *testing.T) {
	testMoveCount(t, InitialiseState(), 119060324, 6)
}

// TODO: Add further positions
// Reference:
// https://www.chessprogramming.org/Perft_Results

func TestMoveCountPosition5Depth1(t *testing.T) {
	testMoveCount(t, getPosition5(), 44, 1)
}

func TestMoveCountPosition5Depth2(t *testing.T) {
	testMoveCount(t, getPosition5(), 1486, 2)
}

func TestMoveCountPosition5Depth3(t *testing.T) {
	testMoveCount(t, getPosition5(), 62379, 3)
}

func TestMoveCountPosition5Depth4(t *testing.T) {
	testMoveCount(t, getPosition5(), 2103487, 4)
}

func TestMoveCountPosition5Depth5(t *testing.T) {
	testMoveCount(t, getPosition5(), 89941194, 5)
}

func testMoveCount(t *testing.T, state State, expected, depth int) {
	moveCount, err := getMoveCount(state, depth)
	if err != nil {
		t.Fatalf("move count failed: %s", err.Error())
	}

	if moveCount != expected {
		t.Fatalf("move count was incorrect:\nexpected=%d\nactual=%d", expected, moveCount)
	}
}

func getMoveCount(state State, depth int) (int, error) {
	if depth == 0 {
		return 1, nil
	}
	depth--

	possibleMoves, err := state.GenerateAllMoves()
	if err != nil {
		return 0, err
	}

	var total int = 0
	for _, newMove := range possibleMoves {
		// Save previous state information which is not contained within a move
		castlingRights := state.CastlingRights
		enPassantSquare := state.EnPassantPosition

		state.DoMove(newMove)
		moveCount, err := getMoveCount(state, depth)
		state.UndoMove(newMove, castlingRights, enPassantSquare)

		if err != nil {
			return 0, err
		}
		total += moveCount
	}

	return total, nil
}

func getPosition5() State {
	return State{
		Board{
			{BlackRook, BlackKnight, BlackBishop, BlackQueen, EmptySquare, BlackKing, EmptySquare, BlackRook},
			{BlackPawn, BlackPawn, EmptySquare, WhitePawn, BlackBishop, BlackPawn, BlackPawn, BlackPawn},
			{EmptySquare, EmptySquare, BlackPawn, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
			{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
			{EmptySquare, EmptySquare, WhiteBishop, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
			{EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare, EmptySquare},
			{WhitePawn, WhitePawn, WhitePawn, EmptySquare, WhiteKnight, BlackKnight, WhitePawn, WhitePawn},
			{WhiteRook, WhiteKnight, WhiteBishop, WhiteQueen, WhiteKing, EmptySquare, EmptySquare, WhiteRook},
		},
		CastlingRights{
			true,
			true,
			false,
			false,
		},
		White,
		PositionOpt{Ok: false},
	}
}
