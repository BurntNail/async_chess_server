package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host   = "db"
	port   = 5432
	user   = "postgres"
	dbname = "async_chess"
)

var GlobalDb *sql.DB
var GlobalDbMutex sync.Mutex

var ClientGetValidCaches map[string]JVec[int] = make(map[string]JVec[int])
var ClientGetCacheMutex sync.RWMutex

func main() {
	loadEnv()
	loadDB()
	defer func() {
		fmt.Println("Exiting app")
		GlobalDb.Close()
	}()
	startServer()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}
func loadDB() {
	GlobalDbMutex.Lock()

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, os.Getenv("DB_PASSWORD"), dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}
	if err := CreateTable(db); err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB!")

	GlobalDb = db
	GlobalDbMutex.Unlock()
}
func startServer() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "all OK"})
	})
	router.GET("/games/:id", APIGetPieces)
	router.POST("/newgame", APINewGame)
	router.POST("/deletegame", APIDeleteGame)
	router.POST("/movepiece", APIMovePiece)
	router.POST("/invalidate", APIInvalidateClientIPCache)

	//TODO: Make this able to run on any server
	if err := router.Run(":12345"); err != nil {
		fmt.Println("Error with Running server:", err.Error())
	}
}
