package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       string `json:"id"`
	Email    string `binding:"required" json:"email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (id, email, password) VALUES (?, ?, ?)`
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
	u.ID = utils.GenerateUUID()
	_, err = stmt.Exec(u.ID, u.Email, u.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByID(id string) (*User, error) {
	query := `SELECT id, email, role FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, id)
	if row == nil {
		return nil, nil
	}

	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, email, role FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Role)
		if err != nil {
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) ValidateCredentials() (*string, error) {
	query := `SELECT id, password, role FROM users WHERE email = ?`
	row := db.DB.QueryRow(query, u.Email, u.Role)
	if row == nil {
		return nil, nil
	}

	var hashedPassword string
	err := row.Scan(&u.ID, &hashedPassword, &u.Role)
	if err != nil {
		return nil, err
	}

	isPasswordValid := utils.CheckPasswordHash(u.Password, hashedPassword)
	if !isPasswordValid {
		return nil, errors.New("invalid credentials")
	}
	return &u.Role, nil
}

func (u *User) Update() error {
	query := `UPDATE users SET email = ?, role = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Email, u.Role, u.ID)
	if err != nil {
		return err
	}

	return nil
}
