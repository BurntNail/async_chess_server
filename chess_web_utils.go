package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func APIGetPieces(c *gin.Context) {
	GlobalDbMutex.Lock()
	defer GlobalDbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if valid_ids, ok := GlobalDbValidCaches[c.ClientIP()]; !ok || valid_ids.IndexOf(id) == -1 {
		//Refresh cache as client ip not present, or id not present

		pieces, err := GetBoard(id, GlobalDb, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, pieces)
		}

		valid_ids.Add(id)
		GlobalDbValidCaches[c.ClientIP()] = valid_ids

		return
	}

	//The client already has a valid cache
	c.JSON(http.StatusAlreadyReported, "")
}

func APINewGame(c *gin.Context) {
	GlobalDbMutex.Lock()
	defer GlobalDbMutex.Unlock()

	var id int
	if err := c.BindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := NewGame(id, GlobalDb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "board successfully created"})

		current := GlobalDbValidCaches[c.ClientIP()]
		current.Remove(id)
		GlobalDbValidCaches[c.ClientIP()] = current
	}
}

// func APIDeleteTable(c *gin.Context) {
// 	DbMutex.Lock()
// 	defer DbMutex.Unlock()

// 	if num, err := DeleteTable(GlobalDb); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
// 	} else {
// 		c.JSON(http.StatusOK, gin.H{"Rows Affected": num})
// 	}
// }

func APIDeleteGame(c *gin.Context) {
	GlobalDbMutex.Lock()
	defer GlobalDbMutex.Unlock()

	var id int
	if err := c.BindJSON(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if num, err := DeleteGame(GlobalDb, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Rows Affected": num})

		current := GlobalDbValidCaches[c.ClientIP()]
		current.Remove(id)
		GlobalDbValidCaches[c.ClientIP()] = current
	}
}

type PieceMove struct {
	ID   int `json:"id"`
	X    int `json:"x"`
	Y    int `json:"y"`
	NewX int `json:"nx"`
	NewY int `json:"ny"`
}

const defaultPMField = -42069

func APIMovePiece(c *gin.Context) {
	GlobalDbMutex.Lock()
	defer GlobalDbMutex.Unlock()

	var move PieceMove = PieceMove{ID: defaultPMField, X: defaultPMField, Y: defaultPMField, NewX: defaultPMField, NewY: defaultPMField}
	if err := c.BindJSON(&move); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	{
		errors := make([]string, 0, 6)

		if move.X == defaultPMField {
			errors = append(errors, "x not set")
			move.X = 0
		} else if move.X < 0 {
			errors = append(errors, "x < 0")
		} else if move.X >= 8 {
			errors = append(errors, "x > 7")
		}
		if move.Y == defaultPMField {
			errors = append(errors, "y not set")
			move.Y = 0
		} else if move.Y < 0 {
			errors = append(errors, "y < 0")
		} else if move.Y >= 8 {
			errors = append(errors, "y > 7")
		}
		if move.NewX == defaultPMField {
			errors = append(errors, "nx not set")
			move.NewX = 0
		} else if move.NewX < 0 {
			errors = append(errors, "nx < 0")
		} else if move.NewX >= 8 {
			errors = append(errors, "nx > 7")
		}
		if move.NewY == defaultPMField {
			errors = append(errors, "ny not set")
			move.NewY = 0 //must set these to 0 to avoid user confusion surrounding the weird default value I set to check for this
		} else if move.NewY < 0 {
			errors = append(errors, "ny < 0")
		} else if move.NewY >= 8 {
			errors = append(errors, "ny > 7")
		}

		if len(errors) == 0 {
			if move.X == move.NewX && move.Y == move.NewY {
				errors = append(errors, "new position same as old position")
			}
		}

		if move.ID == defaultPMField {
			errors = append(errors, "id not set")
			move.ID = 0
		} else if move.ID < 0 {
			errors = append(errors, "id < 0")
		}

		if len(errors) != 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error validating fields": errors, "user input parsed": move})
			return
		}
	}

	if wasTaken, err := MovePiece(GlobalDb, move.ID, move.X, move.Y, move.NewX, move.NewY); err != nil {
		var statusCode int
		if err.Error() == "unable to find piece in given position" {
			statusCode = http.StatusBadRequest
		} else {
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
	} else {
		if wasTaken {
			c.JSON(http.StatusOK, "Piece was taken")
		} else {
			c.JSON(http.StatusOK, "Piece was not taken")
		}

		current := GlobalDbValidCaches[c.ClientIP()]
		current.Remove(move.ID)
		GlobalDbValidCaches[c.ClientIP()] = current
	}
}
