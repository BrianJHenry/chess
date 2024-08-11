package main

import (
	"fmt"

	"github.com/BrianJHenry/chess/internal/chess"
	"github.com/BrianJHenry/chess/internal/engine"
)

func main() {
	fmt.Println("Type Q or q to quit. ")

	isUseEngine, quit := resolveUseEngine()

	if quit {
		return
	}

	userColor, quit := resolveUserColor()

	if quit {
		return
	}

	var engineVersion string
	if isUseEngine {
		engineVersion, quit = resolveEngine()
		if quit {
			return
		}
	}

	engine, err := engine.ResolveEngine(engineVersion)
	if err != nil {
		fmt.Printf("Error resolving engine %s\n", engineVersion)
		return
	}

	game := chess.InitialiseGame()
	for !game.Checkmate && !game.Stalemate {
		fmt.Println(chess.BoardToDisplayString(game.State.Board))

		var move chess.Move
		if game.State.ActiveColor == userColor {
			move, quit = getUserMove(game.State, game.PossibleMoves)
			if quit {
				return
			}

		} else {
			move = engine.ChooseMove(game.State, game.PossibleMoves)
		}

		err = game.DoMove(move)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if game.Stalemate {
		fmt.Println("Stalemate")
	} else if game.Checkmate {
		if game.State.ActiveColor == chess.White {
			fmt.Println("Black wins!")
		} else {
			fmt.Println("White wins!")
		}
	}

	fmt.Scan()
}

func getUserMove(state chess.State, possibleMoves []chess.Move) (move chess.Move, quit bool) {
	for {
		var userMove string
		fmt.Print("Input move: ")
		fmt.Scanln(&userMove)

		if userMove == "q" || userMove == "Q" {
			return chess.Move{}, true
		}

		move, err := chess.AlgebraicNotation(userMove).ToMove(state)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		isMoveLegal := false
		for _, legalMove := range possibleMoves {
			if move == legalMove {
				isMoveLegal = true
			}
		}

		if !isMoveLegal {
			fmt.Println("Illegal move.")
			continue
		}

		return move, false
	}
}

func resolveUserColor() (color chess.Color, quit bool) {
	for {
		var userColor string
		fmt.Print("Select color: (white: w; black: b) ")
		fmt.Scanln(&userColor)

		if userColor == "w" || userColor == "W" {
			return chess.White, false
		} else if userColor == "b" || userColor == "B" {
			return chess.Black, false
		} else if userColor == "q" || userColor == "Q" {
			return chess.White, true
		}
	}
}

func resolveUseEngine() (isUseEngine bool, quit bool) {

	for {
		var useEngine string
		fmt.Print("Play engine? (y/n) ")
		fmt.Scanln(&useEngine)

		if useEngine == "y" || useEngine == "Y" {
			return true, false
		} else if useEngine == "n" || useEngine == "N" {
			return false, false
		} else if useEngine == "q" || useEngine == "Q" {
			return false, true
		}
	}
}

func resolveEngine() (engineVersion string, quit bool) {
	for {
		fmt.Print("Select engine version: (o for options) ")
		fmt.Scanln(&engineVersion)

		if engineVersion == "o" {
			fmt.Println("Available versions")
			for _, version := range engine.Versions {
				fmt.Println(version.Name)
			}
		} else if engineVersion == "q" || engineVersion == "Q" {
			return "", true
		}

		invalidVersion := true
		for _, validVersion := range engine.Versions {
			if engineVersion == validVersion.Name {
				invalidVersion = false
			}
		}

		if invalidVersion {
			fmt.Printf("invalid version %s\n", engineVersion)
		} else {
			break
		}
	}

	return engineVersion, false
}
