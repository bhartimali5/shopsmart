package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

// these are the fields for product: name, description, price, category, stock quantity

type Product struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float32 `json:"price"`
	CategoryID    string  `json:"category_id"`
	StockQuantity int     `json:"stock_quantity"`
}

var products = []Product{}

func (p *Product) Save() error {

	query := `INSERT INTO products (id, name, description, price, category_id, stock_quantity) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	p.ID = utils.GenerateUUID()
	_, err = stmt.Exec(p.ID, p.Name, p.Description, p.Price, p.CategoryID, p.StockQuantity)
	if err != nil {
		return err
	}

	return nil

}

func GetAllproducts() ([]Product, error) {
	query := `SELECT id, name, description, price, category_id, stock_quantity FROM products`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CategoryID, &product.StockQuantity)
		if err != nil {
			continue
		}
		products = append(products, product)
	}
	return products, nil
}

func GetProductByID(id string) (*Product, error) {
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

func (p *Product) Update() error {
	query := `UPDATE products SET name = ?, description = ?,  price = ?, category_id = ?, stock_quantity = ? WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Name, p.Description, p.Price, p.CategoryID, p.StockQuantity, p.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Delete() error {
	query := `DELETE FROM products WHERE id = ?`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetProductsByCategoryID(categoryID string) ([]Product, error) {
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
