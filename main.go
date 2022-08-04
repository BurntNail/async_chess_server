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
var DbMutex sync.Mutex

func main() {
	loadEnv()
	loadDB()
	startServer()
}

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}
func loadDB() {
	DbMutex.Lock()

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
	DbMutex.Unlock()
}
func startServer() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.GET("/games/:id", APIGetPieces)
	router.GET("/newgame/:id", APINewGame)
	router.GET("/deletegame/:id", APIDeleteGame)
	router.GET("/deleteall", APIDeleteTable)

	router.Run("109.74.205.63:12345")

}
