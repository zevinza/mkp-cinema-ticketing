package transaction

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
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
	id := lib.StringToUUID(c.Params("id"))

	api := new(model.TransactionPaymentAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	trx := model.Transaction{}
	db.Where(`id = ?`, id).Take(&trx)
	if lib.RevStr(trx.TransactionStatus) == "paid" || lib.RevStr(trx.TransactionStatus) == "cancelled" {
		return lib.ErrorNotAllowed(c)
	}

	var data model.TransactionPayment
	lib.Merge(api, &data)
	data.TransactionID = id
	if data.PaidAt == nil {
		data.PaidAt = lib.StrfmtNow()
	}
	if data.ReferenceNo == nil {
		data.ReferenceNo = model.GenRefCount("Payment", db)
	}

	if err := db.Create(&data).Error; err != nil {
		return lib.ErrorConflict(c, err)
	}

	amount := float64(0)
	db.Model(&model.TransactionPayment{}).Select(`sum(paid_amount) as amount`).Where(`transaction_id = ?`, id).Find(&amount)

	if amount >= lib.RevFloat64(trx.TotalPrice) {
		trx.TransactionStatus = lib.Strptr("paid")
	}
	trx.TotalPaid = lib.Float64ptr(amount)

	db.Updates(&trx)
	data.Transaction = &trx

	return lib.OK(c, data)
}
