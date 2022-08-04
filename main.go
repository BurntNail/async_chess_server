package main

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "localdb"
	dbname   = "async_chess"
)

var GlobalDb *sql.DB
var DbMutex sync.Mutex

func main() {
	DbMutex.Lock()
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}
	if err := CreateTable(db); err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB!")

	GlobalDb = db
	DbMutex.Unlock()

	router := gin.Default()
	router.GET("/games/:id", APIGetPieces)
	router.GET("/newgame/:id", APINewGame)

	router.Run("109.74.205.63:12345")
}
