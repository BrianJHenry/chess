package chess

type Turn bool

// TODO: Does this actually make sense
const (
	White Turn = false
	Black Turn = true
)

type State struct {
	Board        Board
	WhiteCastled bool
	BlackCastled bool
	Turn         bool
	Previous     EncodedMove
}
