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

	db.SetConnMaxLifetime(1)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)
	return db
}
