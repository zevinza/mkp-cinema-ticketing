package model

import "api/app/lib"

type Movie struct {
	Base
	DataOwner
	MovieAPI
}

type MovieAPI struct {
	Name     *string  `json:"name,omitempty"`
	Code     *string  `json:"code,omitempty" gorm:"index:movie_code_unique,unique,where:deleted_at is null"`
	Duration *float64 `json:"duration,omitempty"`
	Category *string  `json:"category,omitempty"`
	Producer *string  `json:"producer,omitempty"`
	Director *string  `json:"director,omitempty"`
}

func (b *Movie) Seed() *[]Movie {
	res := []Movie{
		{
			Base: Base{
				ID: ParseStrToUUID(movie1),
			},
			MovieAPI: MovieAPI{
				Name:     lib.Strptr("Avatar 2"),
				Code:     lib.Strptr("avatar_2"),
				Duration: lib.Float64ptr(187),
				Category: lib.Strptr("R"),
				Producer: lib.Strptr("James Cameron"),
				Director: lib.Strptr("James Cameron"),
			},
		},
		{
			Base: Base{
				ID: ParseStrToUUID(movie2),
			},
			MovieAPI: MovieAPI{
				Name:     lib.Strptr("Luck"),
				Code:     lib.Strptr("luck"),
				Duration: lib.Float64ptr(176),
				Category: lib.Strptr("S"),
				Producer: lib.Strptr("Disney"),
				Director: lib.Strptr("Sam Raimi"),
			},
		},
	}
	return &res
}
