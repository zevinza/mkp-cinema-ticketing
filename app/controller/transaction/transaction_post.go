package transaction

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// PostTransactionPayment godoc
// @Summary Post new Transaction
// @Description Post new Transaction
// @Param cart_id path string true "Cart ID"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.Transaction "Transaction data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /transactions/{cart_id} [post]
// @Tags Transaction
func PostTransaction(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	rdb := services.REDIS
	id := c.Params("id")
	trxID := lib.GenUUID()
	userID := lib.GetXUserID(c)

	if _, err := rdb.Get(context.Background(), "cart_"+id).Result(); err != nil {
		return lib.ErrorHTTP(c, 410, "your session is expired")
	}

	cart := model.Cart{}
	if row := db.Where(`id = ?`, id).
		Take(&cart); row.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	show := model.ShowSchedule{}
	db.Where(`id = ?`, cart.ShowScheduleID).
		Preload("Movie").
		Preload("CinemaLocation").
		Preload("Theater").
		Take(&show)

	seatIDs := cart.GetSeat()
	seats := []model.Seat{}
	db.Where(`id IN ?`, seatIDs).Find(&seats)

	user := model.User{}
	db.Where(`id = ?`, userID).Take(&user)

	transaction := model.Transaction{
		Base: model.Base{
			ID: trxID,
		},
		DataOwner: model.DataOwner{
			CreatorID:  userID,
			ModifierID: userID,
		},
		TransactionAPI: model.TransactionAPI{
			UserID:            userID,
			TransactionStatus: lib.Strptr("Unpaid"),
			TotalTicketPrice:  cart.TotalPrice,
			TotalDiscount:     lib.Float64ptr(0),
			TotalFee:          lib.Float64ptr(float64(viper.GetInt("ADMIN_FEE"))),
			TotalPaid:         lib.Float64ptr(0),
			ContactName:       user.FirstName,
			ContactDetail:     user.PhoneNumber,
		},
	}
	transaction.TotalPrice = lib.Float64ptr(lib.RevFloat64(transaction.TotalTicketPrice) - lib.RevFloat64(transaction.TotalDiscount) + lib.RevFloat64(transaction.TotalFee))

	err := db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		tickets := []model.Ticket{}
		for _, seat := range seats {
			tickets = append(tickets, model.Ticket{
				TicketAPI: model.TicketAPI{
					TransactionID:  transaction.ID,
					ShowScheduleID: cart.ShowScheduleID,
					IsPrinted:      lib.Boolptr(false),
					SeatCode:       seat.Code,
					MovieName:      show.Movie.Name,
					LocationName:   show.CinemaLocation.Name,
					TheaterName:    show.Theater.Name,
					Price:          show.Price,
					Date:           show.Date,
					StartTime:      show.StartTime,
				},
			})
		}
		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return lib.ErrorConflict(c, err)
	}

	return lib.OK(c)
}
