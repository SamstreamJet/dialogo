package connector

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var deadDBError error = errors.New("db connection is corrupted")
var internalDBError error = errors.New("internal db error")

var db *sqlx.DB = nil

func Connect() error {
	dbUser := os.Getenv("dbUser")
	dbHost := os.Getenv("dbHost")
	dbPass := os.Getenv("dbPass")
	dbName := os.Getenv("dbName")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", dbUser, dbName, dbPass, dbHost)

	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return deadDBError
	}

	return nil
}

func Alive() error {
	if db == nil {
		return deadDBError
	}
	err := db.Ping()
	if err != nil {
		return deadDBError
	}
	return nil
}

func Query(query string, args ...any) (*sqlx.Rows, error) {
	if db == nil {
		return nil, deadDBError
	}
	fmt.Printf("Attempt to execute query: %s --- ", query)
	result, err := db.Queryx(query, args...)
	if err != nil {
		fmt.Print("ERROR\n")
		return nil, internalDBError
	}
	fmt.Print("SUCCESS\n")
	return result, nil
}

func Select(dest interface{}, query string, args ...any) (error) {
	if db == nil {
		return deadDBError
	}
	fmt.Printf("Attempt to execute select query: %s --- ", query)
	err := db.Select(dest, query, args...)
	if err != nil {
		fmt.Print("ERROR\n")
		fmt.Println(err.Error())
		return internalDBError
	}
	fmt.Print("SUCCESS\n")
	return nil
}

func GetConnx(ctx context.Context) (*sqlx.Conn, error) {
	conn, err := db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// func Select(dest interface{}, query string) error {
// 	if db == nil {
// 		return deadDBError
// 	}
// 	err := db.Select(&dest, query)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
