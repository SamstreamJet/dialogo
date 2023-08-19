package connector

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect() {
    dbUser := os.Getenv("dbUser")
    dbHost := os.Getenv("dbHost")
    dbPass := os.Getenv("dbPass")
    dbName := os.Getenv("dbName")

    connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", dbUser, dbName, dbPass, dbHost)

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
}

func Ping() {
    err := db.Ping()
    if err != nil {
        log.Fatal(err)
    }
}
