package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connected successfully!")

	DB = db
}
