package main

import (
	"errors"
	"fmt"
	"math"
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

// YPOS needs to be from 0 to 7 inclusive
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

		for x := 0; x < 8; x++ {
			board = append(board, Piece{
				Y:       1,
				X:       x,
				Kind:    pawnName,
				IsWhite: false,
			})
			board = append(board, Piece{
				Y:       6,
				X:       x,
				Kind:    pawnName,
				IsWhite: true,
			})
		}
	}

	for x := 0; x < 8; x++ {
		kind, err := IndexToName(StartYPosToIndex(x))
		if err != nil {
			return nil, err
		}

		board = append(board, Piece{
			Y:       0,
			X:       x,
			Kind:    kind,
			IsWhite: false,
		})
		board = append(board, Piece{
			Y:       7,
			X:       x,
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

type SQLPiece struct {
	x         int
	y         int
	kind      int
	is_white  bool
	parent_id int
}

type Board []Piece

func CheckValidMoveNonPawn(current SQLPiece, newX, newY int) bool {
	if current.kind == PAWN {
		fmt.Printf("Pawn passed to CheckValidMoveNotPawn: %v", current)

		return false
	}

	rx := current.x - newX
	ry := current.y - newY
	dx := AbsInt(rx)
	dy := AbsInt(ry)

	// moveX := 1
	// moveY := 1
	// if rx < 0 {
	// 	moveX = -1
	// }
	// if ry < 0 {
	// 	moveY = -1
	// }
	//TODO: not everything can castle

	bishop := dx == dy
	rook := (dx != 0 && dy == 0) || (dx == 0 && dy != 0)
	queen := bishop || rook
	switch current.kind {
	case BISHOP:
		return bishop
	case KNIGHT:
		return (dx == 2 && dy == 1) || (dx == 1 && dy == 2)
	case ROOK:
		return rook
	case QUEEN:
		return queen
	case KING:
		return queen && (PythagDist(dx, dy) < math.Sqrt2)
	default:
	}

	return false
}
func CheckValidMovePawn(current SQLPiece, newX, newY int, takesPiece bool) bool {
	if current.kind != PAWN {
		if kind, err := IndexToName(current.kind); err != nil {
			fmt.Printf("Non-Valid type passed to CheckValidMoveNotPawn: %v", current)
		} else {
			fmt.Printf("%v passed to CheckValidMoveNotPawn: %v", kind, current)
		}
		return false
	}

	maxYDst := 1
	if (current.is_white && current.y == 6) || (!current.is_white && current.y == 1) {
		maxYDst = 2
	}
	dstMovedY := current.y - newY
	if !current.is_white {
		dstMovedY *= -1
	}

	if dstMovedY < 1 || dstMovedY > maxYDst {
		return false
	}

	if takesPiece {
		xdst := AbsInt(current.x - newX)
		return xdst == 1
	} else {
		return current.x == newX
	}
}

func CheckValidMove(current SQLPiece, nx, ny int, takesPiece bool) bool {
	if current.kind == PAWN {
		return CheckValidMovePawn(current, nx, ny, takesPiece)
	} else {
		return CheckValidMoveNonPawn(current, nx, ny)
	}
}
