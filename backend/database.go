package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// db es la instancia global de conexión a la base de datos (*sql.DB es thread-safe).
var db *sql.DB

// connectDB inicializa la conexión con la base de datos PostgreSQL.
// Lee los parámetros de conexión de las variables de entorno.
// Intenta reconectar varias veces si la conexión inicial falla (útil con Docker Compose).
func connectDB() error {
	var err error

	// Retrieve connection details from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Basic check if environment variables are set
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	// If all retries fail, return the last error
	fmt.Println("Conexión a PostgreSQL exitosa")
	return nil
}

// closeDB cierra la conexión a la base de datos si está abierta.
func closeDB() {
	if db != nil {
		db.Close()
		fmt.Println("Conexión a la base de datos cerrada")
	}
}