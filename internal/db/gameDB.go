package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *sql.DB

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func connect() error {
	host := "localhost"
	port := 5432
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASS")
	dbname := os.Getenv("PGDB")

	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	//воткнуть потом при подключении
	//defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}
