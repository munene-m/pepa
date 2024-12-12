package models

import "time"

type Invoice struct {
	BillTo string `json:"bill_to"`
	BillFrom string `json:"bill_from"`
	InvoiceNumber string `json:"invoice_number"`
	InvoiceDate string `json:"invoice_date"`
	DueDate string `json:"due_date"`
	InvoiceSubTotal float64 `json:"invoice_sub_total"`
	Tax float64 `json:"tax"`
	InvoiceTotal float64 `json:"invoice_total"`
	InvoiceStatus string `json:"invoice_status"`
	Notes string `json:"notes"`
	Items []Item `json:"items"`
	Currency string `json:"currency"`
	BillToLogo string `json:"bill_to_logo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Item struct {
	ItemDescription string `json:"item_description"`
	ItemQuantity int `json:"item_quantity"`
	ItemRate float64 `json:"item_rate"`
	Amount float64 `json:"amount"`
}