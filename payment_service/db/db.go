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
	// create payments table query
	createPaymentTableQuery := `CREATE TABLE IF NOT EXISTS PAYMENTS(
	"id" TEXT PRIMARY KEY NOT NULL,
	"order_id" TEXT NOT NULL,
	"user_id" TEXT NOT NULL,
	"cart_id" TEXT NOT NULL,
	"payment_status" VARCHAR NOT NULL,
	"created_at" TEXT NOT NULL,
	"updated_at" TEXT NOT NULL,
	UNIQUE (order_id, user_id)

	);`
	ExecQuery(createPaymentTableQuery)
}
