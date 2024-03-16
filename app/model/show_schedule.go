package model

import (
	"api/app/lib"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShowSchedule struct {
	Base
	DataOwner
	ShowScheduleAPI
	Movie          *Movie          `json:"movie,omitempty" gorm:"foreignKey:ID;references:MovieID"`
	CinemaLocation *CinemaLocation `json:"cinema_location,omitempty" gorm:"foreignKey:ID;references:CinemaLocationID"`
	Theater        *Theater        `json:"theater,omitempty" gorm:"foreignKey:ID;references:TheaterID"`
	Seats          *[]Seat         `json:"price,omitempty" gorm:"foreignKey:ShowScheduleID;references:ID"`
}

type ShowScheduleAPI struct {
	Date             *strfmt.Date     `json:"date,omitempty" swaggertype:"string" format:"date"`
	StartTime        *strfmt.DateTime `json:"start_time,omitempty" swaggertype:"string" format:"date-time"`
	MovieID          *uuid.UUID       `json:"movie_id,omitempty"`
	CinemaLocationID *uuid.UUID       `json:"cinema_location_id,omitempty"`
	TheaterID        *uuid.UUID       `json:"theater_id,omitempty"`
	Price            *float64         `json:"price,omitempty"`
	IsStart          *bool            `json:"is_start,omitempty"`
	UniqueCode       *string          `json:"unique_code,omitempty" gorm:"index:show_schedule_code_unique,unique,where:deleted_at is null"`
}

func (b *ShowSchedule) BeforeCreate(tx *gorm.DB) error {
	layouts := []SeatLayout{}
	tx.Find(&layouts)
	b.ID = lib.GenUUID()
	if b.UniqueCode == nil {
		b.UniqueCode = lib.Strptr(lib.RandomChars(20))
	}

	seats := []Seat{}
	for _, layout := range layouts {
		seats = append(seats, Seat{
			SeatAPI: SeatAPI{
				ShowScheduleID: b.ID,
				SeatLayoutID:   layout.ID,
				Code:           layout.Code,
				Row:            layout.Row,
				Column:         layout.Column,
				Characteristic: layout.Characteristic,
				IsAvailable:    lib.Boolptr(true),
			},
		})
	}
	if err := tx.CreateInBatches(&seats, 100).Error; err != nil {
		return err
	}
	return b.Base.BeforeCreate(tx)
}

func (b *ShowSchedule) Seed() *[]ShowSchedule {
	res := []ShowSchedule{
		{
			ShowScheduleAPI: ShowScheduleAPI{
				Date:             lib.Dateptr(strfmt.Date((time.Now()))),
				StartTime:        lib.StrfmtNow(),
				MovieID:          ParseStrToUUID(movie1),
				CinemaLocationID: ParseStrToUUID(location1),
				TheaterID:        ParseStrToUUID(theater1),
				Price:            lib.Float64ptr(35000),
				IsStart:          lib.Boolptr(false),
				UniqueCode:       lib.Strptr("sdfs8d7sdjw93"),
			},
		},
		{
			ShowScheduleAPI: ShowScheduleAPI{
				Date:             lib.Dateptr(strfmt.Date((time.Now()))),
				StartTime:        lib.StrfmtNow(),
				MovieID:          ParseStrToUUID(movie2),
				CinemaLocationID: ParseStrToUUID(location2),
				TheaterID:        ParseStrToUUID(theater1),
				Price:            lib.Float64ptr(40000),
				IsStart:          lib.Boolptr(false),
				UniqueCode:       lib.Strptr("sdfhsd7ewel50032j"),
			},
		},
	}
	return &res
}
