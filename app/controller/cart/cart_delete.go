package cart

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"context"

	"github.com/gofiber/fiber/v2"
)

// DeleteCart godoc
// @Summary Delete Cart by id
// @Description Delete Cart by id
// @Param id path string true "Cart ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /carts/{id} [delete]
// @Tags Cart
func DeleteCart(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	rdb := services.REDIS

	var data model.Cart
	result := db.Model(&data).Where("id = ?", c.Params("id")).Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}
	rdb.Del(context.Background(), "cart_"+c.Params("id"))

	seatIDs := data.GetSeat()
	for _, id := range seatIDs {
		rdb.Del(context.Background(), "seat_"+id.String())
	}

	go db.Model(&model.Seat{}).Where(`id IN ?`, seatIDs).UpdateColumn("is_available", "true")

	db.Unscoped().Delete(&data)

	return lib.OK(c)
}
