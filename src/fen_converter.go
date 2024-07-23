package chess

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertFenToState(fenString string) (State, error) {
	// Read and validate tokens
	tokens := strings.Split(fenString, "-")
	if len(tokens) != 6 {
		return State{}, fmt.Errorf("incorrect number of tokens in FEN state: %s", fenString)
	}

	boardRanks := strings.Split(tokens[0], "/")
	if len(boardRanks) != 8 {
		return State{}, fmt.Errorf("incorrect number of ranks in FEN state: %s", fenString)
	}

	color := tokens[1]
	if color != "w" && color != "b" {
		return State{}, fmt.Errorf("invalid active color specified in FEN state: %s", fenString)
	}

	castlingRights := tokens[2]
	if len(castlingRights) > 4 {
		return State{}, fmt.Errorf("too many characters specifying castling rights in FEN state: %s", fenString)
	}

	for _, char := range castlingRights {
		if char != '-' && char != 'K' && char != 'Q' && char != 'k' && char != 'q' {
			return State{}, fmt.Errorf("invalid character specifying castling rights in FEN state: %s", fenString)
		}
	}

	enPassant := tokens[3]
	if len(enPassant) > 2 {
		return State{}, fmt.Errorf("too many characters specifying en passant target square in FEN state: %s", fenString)
	}

	_, err := strconv.Atoi(tokens[4])
	if err != nil {
		return State{}, err
	}

	_, err = strconv.Atoi(tokens[5])
	if err != nil {
		return State{}, err
	}

	// Convert board
	board := Board{}
	for rankIndex, rank := range boardRanks {
		fileIndex := 0
		for _, char := range rank {
			// Handle number of empty squares
			if number, err := strconv.Atoi(string(char)); err == nil {
				for counter := 0; counter < number; counter++ {
					board[fileIndex][rankIndex] = EmptySquare
					fileIndex++
				}
			} else {
				piece, err := convertFenPieceToPiece(char)
				if err != nil {
					return State{}, err
				}

				board[fileIndex][rankIndex] = piece
				fileIndex++
			}
		}
	}

	var turn Turn
	if color == "w" {
		turn = WhiteTurn
	} else if color == "b" {
		turn = BlackTurn
	}

	whiteCanCastleQueenSide := false
	whiteCanCastleKingSide := false
	blackCanCastleQueenSide := false
	blackCanCastleKingSide := false
	for _, char := range castlingRights {
		switch char {
		case 'K':
			whiteCanCastleKingSide = true
		case 'Q':
			whiteCanCastleQueenSide = true
		case 'k':
			blackCanCastleKingSide = true
		case 'q':
			blackCanCastleQueenSide = true
		}
	}

	// TODO: en passant check

	return State{
		board,
		whiteCanCastleKingSide,
		whiteCanCastleQueenSide,
		blackCanCastleKingSide,
		blackCanCastleQueenSide,
		turn,
		0, // TODO: en passant
	}, nil
}

func convertFenPieceToPiece(p rune) (Piece, error) {
	switch p {
	case 'r':
		return BlackRook, nil
	case 'n':
		return BlackKnight, nil
	case 'b':
		return BlackBishop, nil
	case 'q':
		return BlackQueen, nil
	case 'k':
		return BlackKing, nil
	case 'p':
		return BlackPawn, nil
	case 'R':
		return WhiteRook, nil
	case 'N':
		return WhiteKnight, nil
	case 'B':
		return WhiteBishop, nil
	case 'Q':
		return WhiteQueen, nil
	case 'K':
		return WhiteKing, nil
	case 'P':
		return WhitePawn, nil
	}

	return 0, fmt.Errorf("invalid piece character in FEN state: %c", p)
}
