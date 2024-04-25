package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Инициализация соединения Postgres
func (db *DB) Init() error {
	db.name = "Postgres"
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), 
	os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SSLMODE"))

	newDB, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("%s: err connected to database: %v\n", db.name, err)
		return err
	}
	
	log.Printf("%s: Connected to database!", db.name)
	db.sqlxDB = newDB;
	return nil
}
