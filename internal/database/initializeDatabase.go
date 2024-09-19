package database

import (
	"database/sql"
	"fmt"
	"log"
)

func InitializeDataBase() *sql.DB {
	connStr := "postgres://postgres:password@localhost:5432/postgres"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to PostgresSQL!")
	return db
}
