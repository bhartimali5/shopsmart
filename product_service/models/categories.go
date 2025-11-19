package models

import (
	"example.com/rest-api/db"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Category) Save() error {
	query := `INSERT INTO categories (name) VALUES (?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(c.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	c.ID = int(id)
	if err != nil {
		return err
	}

	return nil
}

func GetCategoryByID(id int64) (*Category, error) {
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
