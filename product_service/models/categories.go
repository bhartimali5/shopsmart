package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Category) Save() error {
	query := `INSERT INTO categories (id, name) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	c.ID = utils.GenerateUUID()
	_, err = stmt.Exec(c.ID, c.Name)
	if err != nil {
		return err
	}

	return nil
}

func (c *Category) Delete() error {
	query := `DELETE FROM categories WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetCategoryByID(id string) (*Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ?`
	row := db.DB.QueryRow(query, id)
	if row == nil {
		return nil, nil
	}

	var category Category
	err := row.Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func GetAllCategories() ([]Category, error) {
	query := `SELECT id, name FROM categories`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (c *Category) Update() error {
	query := `UPDATE categories SET name = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Name, c.ID)
	if err != nil {
		return err
	}
	return nil
}
