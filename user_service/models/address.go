package models

import (
	"example.com/rest-api/db"
)

type Address struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	UserID  int    `json:"user_id"`
}

func (a *Address) Save() error {
	query := `INSERT INTO addresses (name, address, user_id) VALUES (?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(a.Name, a.Address, a.UserID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	a.ID = int(id)
	if err != nil {
		return err
	}

	return nil
}

func GetAddressByUserID(userID int64) (*Address, error) {
	query := `SELECT id, name, address, user_id FROM addresses WHERE user_id = ?`
	row, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row == nil {
		return nil, nil
	}
	//An user can have multiple addresses
	var address Address
	for row.Next() {
		err := row.Scan(&address.ID, &address.Name, &address.Address, &address.UserID)
		if err != nil {
			return nil, err
		}
	}
	return &address, nil
}
