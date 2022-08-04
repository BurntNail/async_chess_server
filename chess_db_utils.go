package main

import (
	"database/sql"
	"errors"
	"fmt"
)

func CreateTable(db *sql.DB) error {
	createPieces := `CREATE TABLE IF NOT EXISTS "Pieces" (x int NOT NULL, y int NOT NULL, kind int NOT NULL, is_white bool NOT NULL, parent_id int NOT NULL)`

	_, e := db.Exec(createPieces)
	return e
}

func GetBoard(id int, db *sql.DB, startNewIfNotExists bool) (Board, error) {
	pieces := make([]Piece, 0, 8*4)

	rows, err := db.Query(`SELECT "x", "y", "kind", "is_white" FROM "Pieces" WHERE parent_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var x, y, kind int
		var is_white bool

		if err := rows.Scan(&x, &y, &kind, &is_white); err != nil {
			return nil, err
		}

		// fmt.Println("Found ", x, ", ", y, ", ", IndexToNames(kind), ", ", is_white)
		pieces = append(pieces, Piece{
			X:       x,
			Y:       y,
			Kind:    kind,
			IsWhite: is_white,
		})
	}

	if len(pieces) != 8*4 {
		if startNewIfNotExists {
			fmt.Println("Unable to find enough pieces, new game time!")

			if err := NewGame(id, db); err != nil {
				return nil, err
			} else {
				return GetBoard(id, db, false)
			}
		} else {
			return nil, errors.New("unable to find enough pieces")
		}

	} else {
		return pieces, nil
	}
}

func NewGame(id int, db *sql.DB) error {
	deleteStmt := `DELETE FROM "Pieces" WHERE parent_id = $1`
	if _, err := db.Exec(deleteStmt, id); err != nil {
		return err
	}

	insertStmt := `INSERT INTO "Pieces"("x", "y", "kind", "is_white", "parent_id") VALUES ($1, $2, $3, $4, $5)`
	for x := 0; x < 2; x++ {
		for y := 0; y < 8; y++ {
			if _, e := db.Exec(insertStmt, x, y, StartYPosToIndex(y), false, id); e != nil {
				return e
			}
		}
	}

	for x := 6; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if _, e := db.Exec(insertStmt, x, y, StartYPosToIndex(y), true, id); e != nil {
				return e
			}
		}
	}

	return nil
}
