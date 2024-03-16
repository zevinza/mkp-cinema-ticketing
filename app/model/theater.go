package model

import (
	"api/app/lib"

	"github.com/google/uuid"
)

type Theater struct {
	Base
	DataOwner
	TheaterAPI
}

type TheaterAPI struct {
	CinemaLocationID *uuid.UUID `json:"cinema_location_id,omitempty" gorm:"index:location_theater_code_unique,unique,where:deleted_at is null"`
	Name             *string    `json:"name,omitempty"`
	Code             *string    `json:"code,omitempty" gorm:"index:location_theater_code_unique,unique,where:deleted_at is null"`
}

func (b *Theater) Seed() *[]Theater {
	res := []Theater{
		{
			Base: Base{
				ID: ParseStrToUUID(theater1),
			},
			TheaterAPI: TheaterAPI{
				CinemaLocationID: ParseStrToUUID(location1),
				Name:             lib.Strptr("1"),
				Code:             lib.Strptr("1"),
			},
		},
		{
			Base: Base{
				ID: ParseStrToUUID(theater2),
			},
			TheaterAPI: TheaterAPI{
				CinemaLocationID: ParseStrToUUID(location1),
				Name:             lib.Strptr("2"),
				Code:             lib.Strptr("2"),
			},
		},
		{
			Base: Base{
				ID: ParseStrToUUID(theater1),
			},
			TheaterAPI: TheaterAPI{
				CinemaLocationID: ParseStrToUUID(location2),
				Name:             lib.Strptr("1"),
				Code:             lib.Strptr("1"),
			},
		},
		{
			Base: Base{
				ID: ParseStrToUUID(theater2),
			},
			TheaterAPI: TheaterAPI{
				CinemaLocationID: ParseStrToUUID(location2),
				Name:             lib.Strptr("2"),
				Code:             lib.Strptr("2"),
			},
		},
	}
	return &res
}
