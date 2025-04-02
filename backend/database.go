package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() error {
	var err error

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	fmt.Println("Conexión a PostgreSQL exitosa")
	return nil
}

func closeDB() {
	if db != nil {
		db.Close()
		fmt.Println("Conexión a la base de datos cerrada")
	}
}