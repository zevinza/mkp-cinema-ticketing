package transaction

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PostTransactionCancel godoc
// @Summary Cancel Unpaid Transaction
// @Description Cancel Unpaid Transaction
// @Param id path string true "Transaction ID"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.Transaction "Transaction data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /transactions/{id}/cancel [post]
// @Tags Transaction
func PostTransactionCancel(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.Transaction
	result := db.Model(&data).
		Where(db.Where(model.Transaction{
			Base: model.Base{
				ID: &id,
			},
		})).
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	var countPayment int64
	db.Model(&model.TransactionPayment{}).Where(`transaction_id = ?`, id).Count(&countPayment)
	if countPayment > 0 {
		lib.ErrorNotAllowed(c)
	}

	if err := db.Model(&data).UpdateColumn("transaction_status", "cancelled").Error; err != nil {
		return lib.ErrorConflict(c)
	}

	return lib.OK(c, data)
}
