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
		var x, y int
		var kind_i int
		var is_white bool

		if err := rows.Scan(&x, &y, &kind_i, &is_white); err != nil {
			return nil, err
		}
		kind, err := IndexToName(kind_i)
		if err != nil {
			return nil, err
		}

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
	board, err := DefaultBoard()
	if err != nil {
		return err
	}
	for _, piece := range board {
		index, err := NameToIndex(piece.Kind)
		if err != nil {
			return err //wtf
		}
		if _, e := db.Exec(insertStmt, piece.X, piece.Y, index, piece.IsWhite, id); e != nil {
			return e
		}
	}

	return nil
}

// Can provide a nil function to validate
func getValidPiecesFromRows(rows *sql.Rows, validate func(SQLPiece) bool) ([]SQLPiece, error) {
	slice := make([]SQLPiece, 0)

	var x, y, kind, parent_id int
	var is_white bool
	var sqlp SQLPiece

	for rows.Next() {
		if err := rows.Scan(&x, &y, &kind, &is_white, &parent_id); err != nil {
			return nil, err
		}
		sqlp = SQLPiece{x: x, y: y, kind: kind, is_white: is_white, parent_id: parent_id}

		fmt.Printf("validating %v", sqlp)
		if validate != nil {
			if validate(sqlp) {
				slice = append(slice, sqlp)
			}
		} else {
			slice = append(slice, sqlp)
		}
	}

	return slice, nil
}

// bool signifies whether or not piece was taken
func MovePiece(db *sql.DB, id, x, y, newX, newY int) (bool, error) {
	var pieceTaken bool

	//check whether or not there is a piece to move
	if currentp_rows, err := db.Query(`SELECT * FROM "Pieces" WHERE "x"=$1 AND "y"=$2 AND "parent_id"=$3`, x, y, id); err != nil {
		return false, err
	} else {
		defer currentp_rows.Close()
		if currentpieces, err := getValidPiecesFromRows(currentp_rows, nil); err != nil {
			return false, err
		} else {
			if len(currentpieces) == 0 {
				return false, errors.New("unable to find piece in given position")
			}
			currentpiece := currentpieces[0]
			if currentpiece.x == newX && currentpiece.y == newY {
				return false, errors.New("invalid move")
			}

			if takenp_rows, err := db.Query(`SELECT * FROM "Pieces" WHERE "x"=$1 AND "y"=$2 AND "parent_id"=$3`, newX, newY, id); err != nil {
				return false, err
			} else {
				defer takenp_rows.Close()

				v := func(sqlp SQLPiece) bool {
					return sqlp.is_white != currentpiece.is_white
				}

				if takenpieces, err := getValidPiecesFromRows(takenp_rows, v); err != nil {
					return false, err
				} else {
					pieceTaken = len(takenpieces) == 1

					if !CheckValidMove(currentpiece, newX, newY, pieceTaken) {
						return false, errors.New("invalid move")
					}
				}
			}
		}

	}

	if _, err := db.Exec(`UPDATE "Pieces" SET "x"=-1, "y"=-1 WHERE "x"=$1 AND "y"=$2 AND "parent_id"=$3`, newX, newY, id); err != nil {
		return false, err
	}

	if res, err := db.Exec(`UPDATE "Pieces" SET "x"=$3, "y"=$4 WHERE "x"=$1 AND "y"=$2 AND "parent_id"=$5`, x, y, newX, newY, id); err != nil {
		return pieceTaken, err
	} else {
		if rows, err := res.RowsAffected(); err != nil {
			return pieceTaken, err
		} else {
			fmt.Println("Affected", rows, "rows")
			return pieceTaken, nil
		}
	}
}

func DeleteGame(db *sql.DB, id int) (int, error) {
	return rows_and_error(db, fmt.Sprint(`DELETE FROM "Pieces" WHERE parent_id=`, id))
}

func rows_and_error(db *sql.DB, query string) (int, error) {
	res, err := db.Exec(query)
	if err != nil {
		return -1, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	return int(rows), nil
}
