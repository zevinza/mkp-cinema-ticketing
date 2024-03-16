package model

import (
	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	Base
	DataOwner
	TicketAPI
}

type TicketAPI struct {
	TransactionID   *uuid.UUID       `json:"transaction_id,omitempty"`
	ShowScheduleID  *uuid.UUID       `json:"show_schedule_id,omitempty"`
	ReferenceNumber *string          `json:"reference_number,omitempty"`
	IsPrinted       *bool            `json:"is_printed,omitempty"`
	IsActivated     *bool            `json:"is_activated,omitempty"`
	SeatCode        *string          `json:"seat_code,omitempty"`
	MovieName       *string          `json:"movie_name,omitempty"`
	LocationName    *string          `json:"location_name,omitempty"`
	TheaterName     *string          `json:"theater_name,omitempty"`
	Price           *float64         `json:"price,omitempty"`
	Date            *strfmt.Date     `json:"date,omitempty" swaggertype:"string" format:"date"`
	StartTime       *strfmt.DateTime `json:"start_time,omitempty" swaggertype:"string" format:"date-time"`
}

func (b *Ticket) BeforeCreate(tx *gorm.DB) error {
	if b.ReferenceNumber == nil {
		b.ReferenceNumber = GenRefCount("Ticket", tx)
	}

	return b.Base.BeforeCreate(tx)
}
