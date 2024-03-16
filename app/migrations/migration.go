package migrations

import "api/app/model"

// ModelMigrations models to automigrate
var ModelMigrations = []interface{}{
	model.Cart{},
	model.CinemaLocation{},
	model.Movie{},
	model.ReferenceCount{},
	model.Seat{},
	model.SeatLayout{},
	model.ShowSchedule{},
	model.Theater{},
	model.Ticket{},
	model.Transaction{},
	model.TransactionPayment{},
	model.User{},
}
