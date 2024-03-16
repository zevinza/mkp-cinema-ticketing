package model

import "api/app/lib"

type CinemaLocation struct {
	Base
	DataOwner
	CinemaLocationAPI
}

type CinemaLocationAPI struct {
	Name    *string `json:"name,omitempty"`
	Code    *string `json:"code,omitempty" gorm:"index:cinema_location_code_unique,unique,where:deleted_at is null"`
	Address *string `json:"address,omitempty"`
}

func (b *CinemaLocation) Seed() *[]CinemaLocation {
	res := []CinemaLocation{
		{
			Base: Base{
				ID: ParseStrToUUID(location1),
			},
			CinemaLocationAPI: CinemaLocationAPI{
				Name:    lib.Strptr("Plaza Indonesia"),
				Code:    lib.Strptr("plaza_indonesia"),
				Address: lib.Strptr("Jl. M.H. Thamrin No.28 30, RT.9/RW.5, Gondangdia, Menteng, Central Jakarta City, Jakarta 10350"),
			},
		},
		{
			Base: Base{
				ID: ParseStrToUUID(location2),
			},
			CinemaLocationAPI: CinemaLocationAPI{
				Name:    lib.Strptr("Plaza Surabaya"),
				Code:    lib.Strptr("plaza_surabaya"),
				Address: lib.Strptr("Jl. Pemuda No.31-37 Lt. 5, Embong Kaliasin, Kec. Genteng, Surabaya, Jawa Timur 60271"),
			},
		},
	}
	return &res
}
