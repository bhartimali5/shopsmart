package models

import (
	"example.com/rest-api/db"
)

// these are the fields for product: name, description, price, category, stock quantity

type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Price         string `json:"price"`
	CategoryID    int    `json:"category_id"`
	StockQuantity string `json:"stock_quantity"`
}

var products = []Product{}

func (e *Product) Save() error {

	query := `INSERT INTO products (name, description, price, category_id, stock_quantity) VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(e.Name, e.Description, e.Price, e.CategoryID, e.StockQuantity)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	e.ID = int(id)
	if err != nil {
		return err
	}

	return nil

}

func GetAllproducts() ([]Product, error) {
	query := `SELECT id, name, description, price, category_id, stock_quantity FROM products`
	row, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var product Product
		err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.StockQuantity)
		if err != nil {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProductByID(id int64) (*Product, error) {
	query := `SELECT id, name, description, price, category_id, stock_quantity FROM products WHERE id = ?`
	row := db.DB.QueryRow(query, id)
	if row == nil {
		return nil, nil
	}

	var product Product
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.StockQuantity)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (e Product) Update() error {
	query := `UPDATE products SET name = ?, description = ?,  price = ?, category_id = ?, stock_quantity = ?, WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.Name, e.Description, e.Price, e.CategoryID, e.StockQuantity, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (e Product) Delete() error {
	query := `DELETE FROM products WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetProductsByCategoryID(categoryID int64) ([]Product, error) {
	query := `SELECT id, name, description, price, category_id, stock_quantity FROM products WHERE category_id = ?`
	rows, err := db.DB.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.StockQuantity)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
