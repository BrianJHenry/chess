package chess

/*
Move types
*/

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

const threeBitMask uint8 = 0b00000111
const fourBitMask uint8 = 0b00001111

/*
3 bits for X1
3 bits for Y1
3 bits for X2
3 bits for Y2
4 bits for flag
*/
type EncodedMove int16

type NotatedMove string

type Move struct {
	Start, End Position
	Flag       MoveFlag
}

type MoveWithInfo struct {
	Move     Move
	Moved    Piece
	Captured Piece
}

/*
Functions to translate between move types
*/

func (move Move) EncodeMove() (encoded EncodedMove) {
	encoded = 0
	encoded |= EncodedMove((threeBitMask & uint8(move.Start.X)))
	encoded |= EncodedMove((threeBitMask & uint8(move.Start.Y))) << 3
	encoded |= EncodedMove((threeBitMask & uint8(move.End.X))) << 6
	encoded |= EncodedMove((threeBitMask & uint8(move.End.Y))) << 9
	encoded |= EncodedMove((fourBitMask & uint8(move.Flag))) << 12
	return
}

// TODO
func (move Move) NotateMove(board *Board) (notated NotatedMove) {
	return
}

// TODO
func (notated NotatedMove) EncodeNotatedMove(board *Board) (encoded EncodedMove) {
	return
}

// TODO
func (notated NotatedMove) DenotateMove(board *Board) (move Move) {
	return
}

func (encoded EncodedMove) DecodeMove() (move Move) {
	move = Move{
		Start: Position{X: int8(encoded & EncodedMove(threeBitMask)), Y: int8(encoded & (EncodedMove(threeBitMask) << 3) >> 3)},
		End:   Position{X: int8(encoded & (EncodedMove(threeBitMask) << 6) >> 6), Y: int8(encoded & (EncodedMove(threeBitMask) << 9) >> 9)},
		Flag:  MoveFlag(encoded & (EncodedMove(threeBitMask) << 12) >> 12),
	}
	return
}

// TODO
func (encoded EncodedMove) NotateEncodedMove(board *Board) (notated NotatedMove) {
	return
}
