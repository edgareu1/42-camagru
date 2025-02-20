package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error
	connStr := "host=database port=5432 user=myuser password=mypassword dbname=mydb sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error opening database: ", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		log.Println("Error opening database: ", err)
		return err
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
