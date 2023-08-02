package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	// TODO fill this in directly or through environment variable
	// Build a DSN e.g. postgres://username:password@url.com:5432/dbName
	DB_DSN = "postgres://nitish:2023@localhost:5432/postgres?sslmode=disable"
)

func createDBConnection() {
	// Read environment variables
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	if dbHost != "" {
		// Construct the connection string
		DB_DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
		fmt.Println("db url: ", DB_DSN)
	}
	var err error
	DB, err = sql.Open("postgres", DB_DSN)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
		log.Fatal(err)
	} else {
		err = initializeDatabase(DB)
		if err != nil {
			log.Fatal("Failed to init db: ", err)
			log.Fatal(err)
		}
		fmt.Println("connected")
	}
	// defer DB.Close()
}

func initializeDatabase(db *sql.DB) error {
	// Read the SQL script file
	sqlScript := `CREATE TABLE IF NOT EXISTS players (
		id serial PRIMARY KEY,
		name character varying(15) NOT NULL,
		country character varying(2) NOT NULL,
		score integer NOT NULL
	  );`

	// Execute the SQL script
	_, err := db.Exec(sqlScript)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully.")
	return nil
}
