package model

import "api/app/lib"

type SeatLayout struct {
	Base
	DataOwner
	SeatLayoutAPI
}

type SeatLayoutAPI struct {
	Code           *string `json:"code,omitempty" gorm:"index:seat_layout_code_unique,unique,where:deleted_at is null"`
	Row            *string `json:"row,omitempty"`
	Column         *string `json:"column,omitempty"`
	Characteristic *string `json:"characteristic,omitempty"`
}

func (b *SeatLayout) Seed() *[]SeatLayout {
	res := []SeatLayout{
		{
			SeatLayoutAPI: SeatLayoutAPI{
				Code:           lib.Strptr("A1"),
				Row:            lib.Strptr("A"),
				Column:         lib.Strptr("1"),
				Characteristic: lib.Strptr("standart seat"),
			},
		},
		{
			SeatLayoutAPI: SeatLayoutAPI{
				Code:           lib.Strptr("A2"),
				Row:            lib.Strptr("A"),
				Column:         lib.Strptr("2"),
				Characteristic: lib.Strptr("standart seat"),
			},
		},
		{
			SeatLayoutAPI: SeatLayoutAPI{
				Code:           lib.Strptr("B1"),
				Row:            lib.Strptr("B"),
				Column:         lib.Strptr("1"),
				Characteristic: lib.Strptr("standart seat"),
			},
		},
		{
			SeatLayoutAPI: SeatLayoutAPI{
				Code:           lib.Strptr("B2"),
				Row:            lib.Strptr("B"),
				Column:         lib.Strptr("B"),
				Characteristic: lib.Strptr("standart seat"),
			},
		},
	}

	return &res
}
