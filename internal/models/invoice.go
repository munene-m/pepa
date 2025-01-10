package models

import (
	"time"
)

type Invoice struct {
    ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    BillTo          string    `gorm:"size:255" json:"bill_to"`
    BillFrom        string    `gorm:"size:255" json:"bill_from"`
    InvoiceNumber   int64     `gorm:"unique" json:"invoice_number"`
    InvoiceDate     string    `json:"invoice_date"`
    DueDate         string    `json:"due_date"`
    InvoiceSubTotal float64   `json:"invoice_sub_total"`
    Tax             float64   `json:"tax"`
    InvoiceTotal    float64   `json:"invoice_total"`
    InvoiceStatus   string    `gorm:"size:50" json:"invoice_status"`
    Notes           string    `gorm:"type:text" json:"notes"`
    Items           []Item    `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"items"`
    Currency        string    `gorm:"size:50" json:"currency"`
    BillToLogo      string    `json:"bill_to_logo"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type Item struct {
    ID              uint    `gorm:"primaryKey;autoIncrement" json:"id"`
    InvoiceID       uint    `gorm:"index;not null" json:"invoice_id"`
    ItemDescription string  `gorm:"size:255" json:"item_description"`
    ItemQuantity    int     `json:"item_quantity"`
    ItemRate        float64 `json:"item_rate"`
    Amount          float64 `json:"amount"`
}