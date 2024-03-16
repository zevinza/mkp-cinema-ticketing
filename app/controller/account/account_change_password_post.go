package account

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// PostAccountChangePassword godoc
// @Summary Change password
// @Description Change password
// @Param data body model.ChangePasswordAPI true "Account data"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /accounts/change-password [post]
// @Tags Account
func PostAccountChangePassword(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())

	api := new(model.ChangePasswordAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	password := lib.PasswordEncrypt(lib.RevStr(api.Password), viper.GetString("SALT"), viper.GetString("AES"))

	if row := db.Model(&model.User{}).
		Where(`email = ?`, api.Email).
		UpdateColumn("password", password); row.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	return lib.OK(c)
}
