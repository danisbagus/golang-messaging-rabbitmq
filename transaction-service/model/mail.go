package model

type MailModel struct {
	CustomerEmail string `json:"customer_email"`
	CustomerName  string `json:"customer_name"`
	TransactionID string `json:"transaction_id"`
}
