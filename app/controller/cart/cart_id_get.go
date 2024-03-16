package cart

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetCartID godoc
// @Summary Get a Cart by id
// @Description Get a Cart by id
// @Param id path string true "Cart ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.Cart "Cart data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /carts/{id} [get]
// @Tags Cart
func GetCartID(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.Cart
	result := db.Model(&data).
		Where(db.Where(model.Cart{
			Base: model.Base{
				ID: &id,
			},
		})).
		Preload("ShowSchedule").
		Preload("ShowSchedule.Movie").
		Preload("ShowSchedule.CinemaLocation").
		Preload("ShowSchedule.Theater").
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	return lib.OK(c, data)
}
