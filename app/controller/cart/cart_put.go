package cart

import (
	"api/app/lib"
	"api/app/model"

	"github.com/gofiber/fiber/v2"
)

// PutCart godoc
// @Summary Update Cart by id
// @Description Update Cart by id
// @Param id path string true "Cart ID"
// @Param data body model.CartAPI true "Cart data"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.Cart "Cart data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /carts/{id} [put]
// @Tags Cart
func PutCart(c *fiber.Ctx) error {
	api := new(model.CartPayload)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	// db := services.DB.WithContext(c.UserContext())
	// id, _ := uuid.Parse(c.Params("id"))

	return lib.OK(c)
}
