package postgresql

import (
	"database/sql"
	"log"

	"github.com/brbarme-shop/brbarmex-rating/config"
)

var db *sql.DB

func NewSqlDB(config config.IConfig) *sql.DB {

	var err error
	db, err = sql.Open(config.DabaseDriverName(), config.DatabaseName())
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(3)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	return db
}
