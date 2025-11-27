package main

import (
	"database/sql"
	"fmt"
	"log"

	"example.com/rest-api/utils"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Update with your DB path
	dbPath := "./api.db"

	// Admin details (change as needed)
	email := "admin@gmail.com"
	password := "admin123"
	role := "admin"

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}
	defer db.Close()

	// Insert query
	query := `
		INSERT INTO users (id, email, password, role)
		VALUES (?, ?, ?, ?)
	`

	_, err = db.Exec(query, utils.GenerateUUID(), email, string(hashedPassword), role)
	if err != nil {
		log.Fatal("Error inserting admin user:", err)
	}

	fmt.Println("✅ Admin user created successfully")
}
