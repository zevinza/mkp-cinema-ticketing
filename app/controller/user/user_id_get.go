package user

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetUserID godoc
// @Summary Get a User by id
// @Description Get a User by id
// @Param id path string true "User ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.User "User data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /users/{id} [get]
// @Tags User
func GetUserID(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	id := c.Params("id")

	if !lib.GetXIsAdmin(c) && lib.GetXUserID(c) != nil {
		id = fmt.Sprint(*lib.GetXUserID(c))
	}

	var data model.User
	result := db.Model(&data).
		Where(`id = ? OR username = ?`, id, id).
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	return lib.OK(c, data)
}
