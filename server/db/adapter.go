package connector

import (
	"database/sql"
	"fmt"
	"log"

    _ "github.com/lib/pq"
)

func Connector(dbUser string, dbPass string, dbName string, dbHost string) *sql.DB {
    connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=verify-full", dbUser, dbName, dbPass, dbHost)

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    
    return db
}
