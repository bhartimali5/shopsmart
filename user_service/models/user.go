package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `binding:"required" json:"email"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	res, err := stmt.Exec(u.Email, u.Password)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	u.ID = int(id)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id int64) (*User, error) {
	query := `SELECT id, email, password FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, id)
	if row == nil {
		return nil, nil
	}

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, email, password FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email)
	if row == nil {
		return nil
	}

	var hashedPassword string
	err := row.Scan(&u.ID, &hashedPassword)
	if err != nil {
		return err
	}

	isPasswordValid := utils.CheckPasswordHash(u.Password, hashedPassword)
	if !isPasswordValid {
		return errors.New("invalid credentials")
	}
	return nil
}
