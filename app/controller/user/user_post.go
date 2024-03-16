package user

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// PostUser godoc
// @Summary Create new User
// @Description Create new User
// @Param data body model.UserAPI true "User data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.User "User data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /users [post]
// @Tags User
func PostUser(c *fiber.Ctx) error {
	api := new(model.UserAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	if !lib.GetXIsAdmin(c) {
		return lib.ErrorUnauthorized(c)
	}

	db := services.DB.WithContext(c.UserContext())

	var data model.User
	lib.Merge(api, &data)
	data.CreatorID = lib.GetXUserID(c)
	data.Password = lib.Strptr(lib.PasswordEncrypt("password", viper.GetString("SALT"), viper.GetString("AES")))

	if err := db.Create(&data).Error; nil != err {
		return lib.ErrorConflict(c, err.Error())
	}

	return lib.Created(c, data)
}
