package chess

import (
	"fmt"
	"strconv"
	"strings"
)

func fenToState(fenString string) (State, error) {
	// Read and validate tokens
	tokens := strings.Split(fenString, " ")
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

	// Half moves
	_, err := strconv.Atoi(tokens[4])
	if err != nil {
		return State{}, err
	}

	// Full moves
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
			if n, err := strconv.Atoi(string(char)); err == nil {
				for c := 0; c < n; c++ {
					board[rankIndex][fileIndex] = EmptySquare
					fileIndex++
				}
			} else {
				piece, err := fenPieceToPiece(char)
				if err != nil {
					return State{}, err
				}

				board[rankIndex][fileIndex] = piece
				fileIndex++
			}
		}
	}

	// Active color
	var activeColor Color
	if color == "w" {
		activeColor = White
	} else if color == "b" {
		activeColor = Black
	}

	// Castling rights
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

	// En passant
	var enPassantSquare PositionOpt
	if enPassant == "-" {
		enPassantSquare = PositionOpt{Ok: false}
	} else {
		pos, err := stringToPosition(enPassant)
		if err != nil {
			return State{}, err
		}

		enPassantSquare = PositionOpt{
			Ok:       true,
			Position: pos,
		}
	}

	return State{
		board,
		CastlingRights{
			whiteCanCastleKingSide,
			whiteCanCastleQueenSide,
			blackCanCastleKingSide,
			blackCanCastleQueenSide,
		},
		activeColor,
		enPassantSquare,
	}, nil
}

func fenPieceToPiece(p rune) (Piece, error) {
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
