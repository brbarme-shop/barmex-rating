package postgresql

import (
	"database/sql"
	"log"
)

var db *sql.DB

func NewSqlDB(sourceName string) *sql.DB {

	var err error
	db, err = sql.Open("postgres", sourceName)
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(3)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	return db
}
