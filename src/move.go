package chess

type MoveFlag byte

const (
	None MoveFlag = iota
	EnPassant
	Castle
)

/*
3 bits for x
3 bits for y
2 bits for flag
*/
type EncodedMove int8

type NotatedMove string

type Move struct {
	Start, End Position
	Flag       MoveFlag
}

// TODO
func (move Move) EncodeMove() (encoded EncodedMove, err error) {
	return
}

// TODO
func (move Move) NotateMove() (notated NotatedMove, err error) {
	return
}

// TODO
func (notated NotatedMove) EncodeNotatedMove(board *Board) (encoded EncodedMove, err error) {
	return
}

// TODO
func (notated NotatedMove) DenotateMove(board *Board) (move Move, err error) {
	return
}

// TODO
func (encoded EncodedMove) DecodeMove() (move Move, err error) {
	return
}

// TODO
func (encoded EncodedMove) NotateEncodedMove() (notated NotatedMove, err error) {
	return
}
