package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type dbController struct {
	db *sql.DB
}

func NewController() *dbController {
	db, err := connectDatabase()
	if err != nil {
		log.Fatal("Error making connection with database: ", err)
	}
	return &dbController{
		db: db,
	}
}

func connectDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host= %s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", os.Getenv("POSTGRES_USER_NAME"), os.Getenv("POSTGRES_USER_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return db, nil
}
