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

	//Create category table
	createCategoryTableSQL := `CREATE TABLE IF NOT EXISTS categories (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"name" VARCHAR NOT NULL
	);`
	ExecQuery(createCategoryTableSQL)

	// Create products table
	createProductTableSQL := `CREATE TABLE IF NOT EXISTS products (
        "id" INTEGER PRIMARY KEY AUTOINCREMENT,
        "name" VARCHAR NOT NULL,
        "description" TEXT NOT NULL,
        "price"  DOUBLE NOT NULL,
		"stock_quantity" INTEGER NOT NULL,
        "category_id" INTEGER,
		FOREIGN KEY (category_id) REFERENCES categories(id)
    );`
	ExecQuery(createProductTableSQL)
}
