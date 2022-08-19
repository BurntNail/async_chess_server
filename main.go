package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
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

	router.GET("/games/:id", APIGetPieces)
	router.POST("/newgame", APINewGame)
	router.POST("/deletegame", APIDeleteGame)
	router.POST("/movepiece", APIMovePiece)
	router.POST("/invalidate", APIInvalidateClientIPCache)
	// router.GET("/deleteall", APIDeleteTable) //might re-expose later but not for now

	if err := router.Run("109.74.205.63:12345"); err != nil {
		fmt.Println("Error with Running server:", err.Error())
	}
}
