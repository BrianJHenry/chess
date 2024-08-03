package chess

// MoveFlag is a flag added to a move which handles special cases like castling and promotion.
type MoveFlag uint8

const (
	None MoveFlag = iota
	EnPassant
	KingSideCastle
	QueenSideCastle
	PromoteToQueen
	PromoteToRook
	PromoteToBishop
	PromoteToKnight
)

const bitMask3 uint8 = 0b00000111
const bitMask4 uint8 = 0b00001111

/*
3 bits for X1
3 bits for Y1
3 bits for X2
3 bits for Y2
4 bits for flag
*/
type EncodedMove int16

type Move struct {
	Start, End Position
	Flag       MoveFlag
}

type MoveWithInfo struct {
	Move     Move
	Moved    Piece
	Captured Piece
}

// ToEncoded translates a Move struct to a more compressed encoding.
func (move Move) ToEncoded() (encoded EncodedMove) {
	encoded = 0
	encoded |= EncodedMove((bitMask3 & uint8(move.Start.X)))
	encoded |= EncodedMove((bitMask3 & uint8(move.Start.Y))) << 3
	encoded |= EncodedMove((bitMask3 & uint8(move.End.X))) << 6
	encoded |= EncodedMove((bitMask3 & uint8(move.End.Y))) << 9
	encoded |= EncodedMove((bitMask4 & uint8(move.Flag))) << 12
	return
}

// ToMove translates from a compressed encoding to a Move struct.
func (enc EncodedMove) ToMove() Move {
	return Move{
		Start: Position{X: int8(enc & EncodedMove(bitMask3)), Y: int8(enc & (EncodedMove(bitMask3) << 3) >> 3)},
		End:   Position{X: int8(enc & (EncodedMove(bitMask3) << 6) >> 6), Y: int8(enc & (EncodedMove(bitMask3) << 9) >> 9)},
		Flag:  MoveFlag(enc & (EncodedMove(bitMask3) << 12) >> 12),
	}
}

func MoveTouchesSquare(move Move, position Position) bool {
	return move.Start == position || move.End == position
}
