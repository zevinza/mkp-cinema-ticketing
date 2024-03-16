package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	Base
	DataOwner
	TransactionAPI
	Tickets *[]Ticket `json:"tickets,omitempty" gorm:"foreignKey:TransactionID;references:ID"`
}

type TransactionAPI struct {
	UserID            *uuid.UUID `json:"user_id,omitempty" swaggertype:"string" format:"uuid"`
	CartID            *string    `json:"-"`
	InvoiceNo         *string    `json:"invoice_no,omitempty"`
	TransactionStatus *string    `json:"transaction_status,omitempty"`
	BookingCode       *string    `json:"booking_code,omitempty"`
	TotalTicketPrice  *float64   `json:"total_ticket_price,omitempty"`
	TotalDiscount     *float64   `json:"total_discount,omitempty"`
	TotalFee          *float64   `json:"total_fee,omitempty"`
	TotalPrice        *float64   `json:"total_price,omitempty"`
	ContactName       *string    `json:"contact_name,omitempty"`
	ContactDetail     *string    `json:"contact_detail,omitempty"`
}

func (b *Transaction) BeforeCreate(tx *gorm.DB) error {
	if b.InvoiceNo == nil {
		b.InvoiceNo = GenRefCount("Transaction", tx)
	}

	return b.Base.BeforeCreate(tx)
}
