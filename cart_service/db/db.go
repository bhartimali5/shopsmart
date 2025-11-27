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

	//Create cart table
	createCartTableSQL := `CREATE TABLE IF NOT EXISTS carts (
		"id" TEXT PRIMARY KEY ,
		"user_id" TEXT NOT NULL,
		"created_at" TEXT,
		"updated_at" TEXT
	);`
	ExecQuery(createCartTableSQL)

	// Create cart_item table
	createCartItemTableSQL := `CREATE TABLE IF NOT EXISTS cart_items (
		"id" TEXT PRIMARY KEY,
		"cart_id" TEXT NOT NULL,
		"product_id" TEXT NOT NULL,
		"quantity" INTEGER NOT NULL,
		"price" REAL NOT NULL,
		FOREIGN KEY (cart_id) REFERENCES carts(id)
	);`
	ExecQuery(createCartItemTableSQL)
}
