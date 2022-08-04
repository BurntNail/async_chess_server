package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func APIGetPieces(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pieces, err := GetBoard(id, GlobalDb, false)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
		return
	}

	jsoned, err := json.Marshal(pieces)
	if err != nil {
		panic(err) //not server error, so panic!
	}

	c.JSON(http.StatusOK, jsoned) //make way to json these
}

func APINewGame(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := NewGame(id, GlobalDb); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "board successfully createed"})
	}
}
