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

	// Create order table
	CreateOrderTableQuery := `CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY NOT NULL,
		user_id TEXT NOT NULL,
		order_date TEXT NOT NULL,
		cart_id TEXT NOT NULL,
		status TEXT,
		total_amount REAL NOT NULL
	);`
	ExecQuery(CreateOrderTableQuery)

	//Create user_orders relation table
	CreateUserOrdersTableQuery := `CREATE TABLE IF NOT EXISTS user_orders (
		user_id TEXT NOT NULL,
		order_id TEXT NOT NULL,
		PRIMARY KEY (user_id, order_id)
	);`
	ExecQuery(CreateUserOrdersTableQuery)

}
