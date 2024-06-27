package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = ""
	port     = 6543
	user     = ""
	password = ""
	dbname   = ""
)

var uri = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", user, password, host, port, dbname)

func setup() (*sql.DB, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Init() error {
	db, err := setup()
	if err != nil {
		return err
	}
	DB = db
	return nil
}
