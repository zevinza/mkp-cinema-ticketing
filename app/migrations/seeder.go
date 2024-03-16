package migrations

import "api/app/model"

var (
	location   model.CinemaLocation
	movie      model.Movie
	refCount   model.ReferenceCount
	seatLayout model.SeatLayout
	show       model.ShowSchedule
	user       model.User
)

// DataSeeds data to seeds
func DataSeeds() []interface{} {
	return []interface{}{
		location.Seed(),
		movie.Seed(),
		refCount.Seed(),
		seatLayout.Seed(),
		show.Seed(),
		user.Seed(),
	}
}
