package data

import (
	"database/sql"
	"fmt"
)

// connect used to create a db connection to mysql datadb instance
func connect() (db *sql.DB, err error) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Pa88w0rd"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "db"
	db, err = sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8")
	if err != nil {
		fmt.Printf("[ERROR]  Uh oh! failed to connect to db: %v", err)
		return nil, err
	}
	return db, nil
}
