package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

//Get String DSN string for database
func connSting() string {
	EnvParser()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	port := 5432

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, port, dbUser, dbPass, dbName)
}

//Setup Connection With database
func SetupDatabaseConnection() *sql.DB {
	db, err := sql.Open("postgres", connSting())
	if err != nil {
		fmt.Print(err.Error())
	}
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("connected")
	return db
}
