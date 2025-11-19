package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		panic("Could not connect to database.")
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
}

func ExecQuery(query string) {
	_, err := DB.Exec(query)
	if err != nil {
		panic(err)
	}
}

func CreateTables() {

	// Create users table
	createUsersTabelSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"UserName" VARCHAR(100),
		"email" TEXT NOT NULL UNIQUE,
		"password" TEXT NOT NULL
	);`
	ExecQuery(createUsersTabelSQL)

	//Create address table
	createAddressTableSQL := `CREATE TABLE IF NOT EXISTS addresses (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" VARCHAR(100),
		"address" TEXT,
		"user_id" INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	ExecQuery(createAddressTableSQL)
}
