package model

import (
	"api/app/lib"
	"strings"

	"github.com/google/uuid"
)

// Cart Cart
type Cart struct {
	Base
	DataOwner
	CartAPI
	ShowSchedule *ShowSchedule `json:"show_schedule,omitempty" gorm:"foreignKey:ID;references:ShowScheduleID"`
	Seats        *[]Seat       `json:"seats,omitempty"`
}

// CartAPI Cart API
type CartAPI struct {
	UserID         *uuid.UUID `json:"user_id,omitempty" gorm:"type:varchar(256);index:idx_cart_user_show_schedule_unique,unique,where:deleted_at is null;not null" swaggertype:"string" format:"uuid"`          // UserID
	ShowScheduleID *uuid.UUID `json:"show_schedule_id,omitempty" gorm:"type:varchar(256);index:idx_cart_user_show_schedule_unique,unique,where:deleted_at is null;not null" swaggertype:"string" format:"uuid"` // UserID
	SeatIDs        *string    `json:"seat_ids,omitempty"`
	Price          *float64   `json:"price,omitempty"`
	Quantity       *int       `json:"quantity,omitempty"`
	TotalPrice     *float64   `json:"total_price,omitempty"`
	ExpiresIn      *int64     `json:"expires_in,omitempty"`
}

type CartPayload struct {
	ShowScheduleID *uuid.UUID
	SeatIDs        []uuid.UUID
}

func (b *Cart) GetSeat() []uuid.UUID {
	var ids []uuid.UUID
	if b.SeatIDs != nil {
		seats := strings.Split(lib.RevStr(b.SeatIDs), ", ")
		for _, seat := range seats {
			if uid, err := uuid.Parse(seat); err == nil {
				ids = append(ids, uid)
			}
		}
	}
	return ids
}
