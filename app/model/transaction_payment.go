package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionPayment struct {
	Base
	DataOwner
	TransactionPaymentAPI
	Transaction *Transaction `json:"transaction,omitempty" gorm:"foreignKey:ID;references:TransactionID"`
}

type TransactionPaymentAPI struct {
	TransactionID *uuid.UUID       `json:"transaction_id,omitempty" swaggertype:"string" format:"uuid"`
	PaidAmount    *float64         `json:"paid_amount,omitempty"`
	PaidAt        *strfmt.DateTime `json:"paid_at,omitempty" swaggertype:"string" format:"date-time"`
	Notes         *string          `json:"notes,omitempty"`
	Type          *string          `json:"type,omitempty"`
	Via           *string          `json:"via,omitempty"`
	ReferenceNo   *string          `json:"reference_no,omitempty"`
}

func (b *TransactionPayment) BeforeCreate(tx *gorm.DB) error {
	if b.ReferenceNo == nil {
		b.ReferenceNo = GenRefCount("Payment", tx)
	}

	return b.Base.BeforeCreate(tx)
}
