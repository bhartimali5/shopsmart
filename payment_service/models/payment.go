package models

import (
	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type Payment struct {
	ID            string `json:"id"`
	OrderId       string `json:"order_id"`
	UserId        string `json:"user_id"`
	CartId        string `json:"cart_id"`
	PaymentStatus string `json:"payment_status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func (p *Payment) Save() error {
	query := `INSERT INTO payments (id, order_id, user_id, cart_id, payment_status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	p.ID = utils.GenerateUUID()
	p.CreatedAt, err = utils.GetCurrentTime()
	if err != nil {
		return err
	}
	p.UpdatedAt, err = utils.GetCurrentTime()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.ID, p.OrderId, p.UserId, p.CartId, p.PaymentStatus, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	return nil

}

func GetPaymentsForUserId(userId string) ([]Payment, error) {
	query := `SELECT order_id, user_id, cart_id, payment_status, created_at, updated_at FROM PAYMENTS WHERE user_id = ?`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(&payment.OrderId, &payment.UserId, &payment.CartId, &payment.PaymentStatus, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}
