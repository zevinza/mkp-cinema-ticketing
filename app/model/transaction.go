package model

import (
	"api/app/lib"

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
	InvoiceNo         *string    `json:"invoice_no,omitempty"`
	TransactionStatus *string    `json:"transaction_status,omitempty"`
	BookingCode       *string    `json:"booking_code,omitempty"`
	TotalTicketPrice  *float64   `json:"total_ticket_price,omitempty"`
	TotalDiscount     *float64   `json:"total_discount,omitempty"`
	TotalFee          *float64   `json:"total_fee,omitempty"`
	TotalPrice        *float64   `json:"total_price,omitempty"`
	TotalPaid         *float64   `json:"total_paid,omitempty"`
	ContactName       *string    `json:"contact_name,omitempty"`
	ContactDetail     *string    `json:"contact_detail,omitempty"`
}

func (b *Transaction) BeforeCreate(tx *gorm.DB) error {
	if b.InvoiceNo == nil {
		b.InvoiceNo = GenRefCount("Transaction", tx)
	}
	if b.BookingCode == nil {
		b.BookingCode = lib.Strptr(lib.RandomNumber(5))
	}

	return b.Base.BeforeCreate(tx)
}
