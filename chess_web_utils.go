package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func APIGetPieces(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	pieces, err := GetBoard(id, GlobalDb, false)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, pieces)
}

func APINewGame(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	if err := NewGame(id, GlobalDb); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": err})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "board successfully createed"})
	}
}
