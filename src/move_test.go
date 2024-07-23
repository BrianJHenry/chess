package chess

import "testing"

type testData struct {
	move    Move
	encoded EncodedMove
}

func TestEncodeMove(t *testing.T) {
	for _, test := range retrieveTestData() {
		if test.encoded != test.move.EncodeMove() {
			t.Fatalf("Move: %v does not match encoded move: %b", test.move, test.encoded)
		}
	}
}

func TestDecodeMove(t *testing.T) {
	for _, test := range retrieveTestData() {
		if test.move != test.encoded.DecodeMove() {
			t.Fatalf("Encoded move: %b does not match move: %v", test.encoded, test.move)
		}
	}
}

func retrieveTestData() []testData {
	return []testData{
		{
			move: Move{
				Start: Position{1, 1},
				End:   Position{2, 2},
				Flag:  None,
			},
			encoded: 0b0000010010001001,
		},
		{
			move: Move{
				Start: Position{1, 3},
				End:   Position{1, 5},
				Flag:  None,
			},
			encoded: 0b0000101001011001,
		},
		{
			move: Move{
				Start: Position{4, 7},
				End:   Position{6, 7},
				Flag:  KingSideCastle,
			},
			encoded: 0b0010111110111100,
		},
		{
			move: Move{
				Start: Position{4, 0},
				End:   Position{2, 0},
				Flag:  QueenSideCastle,
			},
			encoded: 0b0011000010000100,
		},
		{
			move: Move{
				Start: Position{4, 3},
				End:   Position{5, 2},
				Flag:  EnPassant,
			},
			encoded: 0b0001010101011100,
		},
		{
			move: Move{
				Start: Position{0, 6},
				End:   Position{0, 7},
				Flag:  PromoteToQueen,
			},
			encoded: 0b0100111000110000,
		},
		{
			move: Move{
				Start: Position{0, 6},
				End:   Position{0, 7},
				Flag:  PromoteToKnight,
			},
			encoded: 0b0111111000110000,
		},
		{
			move: Move{
				Start: Position{0, 6},
				End:   Position{0, 7},
				Flag:  PromoteToBishop,
			},
			encoded: 0b0110111000110000,
		},
		{
			move: Move{
				Start: Position{0, 6},
				End:   Position{0, 7},
				Flag:  PromoteToRook,
			},
			encoded: 0b0101111000110000,
		},
	}
}
