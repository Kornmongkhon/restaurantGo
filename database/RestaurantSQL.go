package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// DB is the database connection pool
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Verify the connection
	if err := DB.Ping(); err != nil {
		log.Fatal("Database is unreachable:", err)
	}

	log.Println("Database connection established")
}
