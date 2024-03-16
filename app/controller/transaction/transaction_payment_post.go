package transaction

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"context"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PostTransactionPayment godoc
// @Summary Pay Transaction
// @Description Pay Transaction
// @Param id path string true "Transaction ID"
// @Param data body model.TransactionPaymentAPI true "Payment data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.TransactionPayment "Transaction data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /transactions/{id}/payment [post]
// @Tags Transaction
func PostTransactionPayment(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	rdb := services.REDIS
	id := lib.StringToUUID(c.Params("id"))
	userID := lib.GetXUserID(c)

	api := new(model.TransactionPaymentAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	trx := model.Transaction{}
	db.Where(`id = ?`, id).Take(&trx)
	if lib.RevStr(trx.TransactionStatus) == "paid" || lib.RevStr(trx.TransactionStatus) == "cancelled" {
		return lib.ErrorNotAllowed(c)
	}

	if val, _ := rdb.Get(context.Background(), "cart_"+lib.RevStr(trx.CartID)).Result(); len(val) == 0 {
		return lib.ErrorHTTP(c, 410, "your session is expired")

		// do Refund if user is paid several amount
	}

	var data model.TransactionPayment
	lib.Merge(api, &data)
	data.TransactionID = id
	if data.PaidAt == nil {
		data.PaidAt = lib.StrfmtNow()
	}
	data.CreatorID = userID

	if lib.RevFloat64(api.PaidAmount) >= lib.RevFloat64(trx.TotalPrice) {
		trx.TransactionStatus = lib.Strptr("paid")
		trx.BookingCode = lib.Strptr(lib.RandomNumber(5))

		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&data).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Ticket{}).Where(`transaction_id = ?`, id).UpdateColumn("is_activated", "true").Error; err != nil {
				return err
			}

			cart := model.Cart{}
			tx.Where(`id = ?`, trx.CartID).Take(&cart)
			seatIDs := cart.GetSeat()
			for _, id := range seatIDs {
				rdb.Del(context.Background(), "seat_"+id.String())
			}

			rdb.Del(context.Background(), "cart_"+lib.RevStr(trx.CartID))

			if err := tx.Where(`id = ?`, trx.CartID).Unscoped().Delete(&model.Cart{}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Seat{}).Where(`id IN ?`, seatIDs).UpdateColumn("is_available", "false").Error; err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return lib.ErrorConflict(c, err)
		}
	}

	db.Updates(&trx)
	data.Transaction = &trx

	return lib.OK(c, data)
}
