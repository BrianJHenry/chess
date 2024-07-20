package chess

import "testing"

type testData struct {
	move    Move
	encoded EncodedMove
}

func TestEncodeMove(t *testing.T) {
	for _, test := range retrieveTestData() {
		if (test.encoded != test.move.EncodeMove(&Board{})) {
			t.Fatalf("Move: %v does not match encoded move: %b", test.move, test.encoded)
		}
	}
}

func TestDecodeMove(t *testing.T) {
	for _, test := range retrieveTestData() {
		if (test.move != test.encoded.DecodeMove(&Board{})) {
			t.Fatalf("Encoded move: %b does not match move: %v", test.encoded, test.move)
		}
	}
}

func retrieveTestData() []testData {
	return []testData{
		{
			move: Move{
				Start: Position{X: 1, Y: 1},
				End:   Position{X: 2, Y: 2},
				Flag:  None,
			},
			encoded: 0b0000010010001001,
		},
		{
			move: Move{
				Start: Position{X: 1, Y: 3},
				End:   Position{X: 1, Y: 5},
				Flag:  None,
			},
			encoded: 0b0000101001011001,
		},
		{
			move: Move{
				Start: Position{X: 7, Y: 3},
				End:   Position{X: 7, Y: 1},
				Flag:  Castle,
			},
			encoded: 0b0010001111011111,
		},
	}
}
