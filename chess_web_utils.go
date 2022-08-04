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
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	pieces, err := GetBoard(id, GlobalDb, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pieces) //make way to json these
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "board successfully createed"})
	}
}

func APIDeleteTable(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	if num, err := DeleteTable(GlobalDb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Rows Affected": num})
	}
}

func APIDeleteGame(c *gin.Context) {
	DbMutex.Lock()
	defer DbMutex.Unlock()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if num, err := DeleteGame(GlobalDb, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Rows Affected": num})
	}
}
