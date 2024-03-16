package model

import "github.com/google/uuid"

type Seat struct {
	Base
	DataOwner
	SeatAPI
	ShowSchedule *ShowSchedule `json:"show_schedule,omitempty" gorm:"foreignKey:ID;references:ShowScheduleID"`
	SeatLayout   *SeatLayout   `json:"seat_layout,omitempty" gorm:"foreignKey:ID;references:SeatLayoutID"`
}

type SeatAPI struct {
	ShowScheduleID *uuid.UUID `json:"show_schedule_id,omitempty"`
	SeatLayoutID   *uuid.UUID `json:"seat_layout_id,omitempty"`
	Code           *string    `json:"code,omitempty"`
	Row            *string    `json:"row,omitempty"`
	Column         *string    `json:"column,omitempty"`
	Characteristic *string    `json:"characteristic,omitempty"`
	IsAvailable    *bool      `json:"is_available,omitempty"`
}
