package main

import (
	"errors"
	"strings"
)

const (
	PAWN = 1 << iota
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

func IndexToNames(index int) []string {
	names := make([]string, 0, 6)

	switch {
	case index&PAWN == PAWN:
		names = append(names, "Pawn")
	case index&BISHOP == BISHOP:
		names = append(names, "Bishop")
	case index&KNIGHT == KNIGHT:
		names = append(names, "Knight")
	case index&ROOK == ROOK:
		names = append(names, "Rook")
	case index&QUEEN == QUEEN:
		names = append(names, "Queen")
	case index&KING == KING:
		names = append(names, "King")
	}

	return names
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

type Piece struct {
	X       int
	Y       int
	Kind    int
	IsWhite bool
}

type Board []Piece
