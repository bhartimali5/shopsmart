package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Address struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  string `json:"user_id"`
}

func (a *Address) Save() error {
	query := `INSERT INTO addresses (id, name, address, user_id) VALUES (?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	a.ID = utils.GenerateUUID()
	_, err = stmt.Exec(a.ID, a.Name, a.Address, a.UserID)
	if err != nil {
		return err
	}

	return nil
}

func GetAddressByUserID(userID string) (*Address, error) {
	query := `SELECT id, name, address, user_id 
              FROM addresses 
              WHERE user_id = ?
              LIMIT 1`

	row := db.DB.QueryRow(query, userID)

	var address Address
	err := row.Scan(&address.ID, &address.Name, &address.Address, &address.UserID)

	if err != nil {
		return nil, err
	}

	return &address, nil
}

func (a *Address) Update() error {
	query := `UPDATE addresses SET name = ?, address = ? WHERE user_id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(a.Name, a.Address, a.UserID)
	if err != nil {
		return err
	}
	return nil
}
