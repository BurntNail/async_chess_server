package main

import (
	"errors"
	"strconv"
	"strings"
)

const (
	PAWN = iota
	BISHOP
	KNIGHT
	ROOK
	QUEEN
	KING
)

func NameToIndex(name string) (int, error) {
	switch strings.ToLower(name) {
	case "pawn":
		return PAWN, nil
	case "bishop":
		return BISHOP, nil
	case "knight":
		return KNIGHT, nil
	case "rook":
		return ROOK, nil
	case "queen":
		return QUEEN, nil
	case "king":
		return KING, nil
	default:
		return -1, errors.New("invalid name")
	}
}

func IndexToName(index int) (string, error) {
	switch index {
	case PAWN:
		return "Pawn", nil
	case BISHOP:
		return "Bishop", nil
	case KNIGHT:
		return "Knight", nil
	case ROOK:
		return "Rook", nil
	case QUEEN:
		return "Queen", nil
	case KING:
		return "King", nil
	default:
		return "", errors.New("invalid index: " + strconv.Itoa(index))
	}
}

//YPOS needs to be from 0 to 7 inclusive
func StartYPosToIndex(ypos int) int {
	switch ypos {
	case 0, 7:
		return ROOK
	case 1, 6:
		return KNIGHT
	case 2, 5:
		return BISHOP
	case 3:
		return QUEEN
	case 4:
		return KING
	default:
		return PAWN
	}
}

func DefaultBoard() (Board, error) {
	board := make([]Piece, 0, 8*4)

	{
		pawnName, err := IndexToName(PAWN)
		if err != nil {
			return nil, err
		}

		for y := 0; y < 8; y++ {
			board = append(board, Piece{
				X:       0,
				Y:       y,
				Kind:    pawnName,
				IsWhite: false,
			})
			board = append(board, Piece{
				X:       7,
				Y:       y,
				Kind:    pawnName,
				IsWhite: true,
			})
		}
	}

	for y := 0; y < 8; y++ {
		kind, err := IndexToName(StartYPosToIndex(y))
		if err != nil {
			return nil, err
		}

		board = append(board, Piece{
			X:       1,
			Y:       y,
			Kind:    kind,
			IsWhite: false,
		})
		board = append(board, Piece{
			X:       6,
			Y:       y,
			Kind:    kind,
			IsWhite: true,
		})

	}

	return board, nil
}

type Piece struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Kind    string `json:"kind"`
	IsWhite bool   `json:"is_white"`
}

type Board []Piece
