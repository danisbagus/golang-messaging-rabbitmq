package model

import "time"

type TransactionModel struct {
	TransactionID   string    `json:"transaction_id"`
	ProductID       int64     `json:"product_id"`
	Quantity        int64     `json:"quantity"`
	TransactionDate time.Time `json:"transaction_date"`
}
