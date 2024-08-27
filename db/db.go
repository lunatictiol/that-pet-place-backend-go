package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewMySqlStorage(connStr string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cant connect to db")
		return nil, err
	}
	return db, nil
}
