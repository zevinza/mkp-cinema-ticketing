package account

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// PostAccountRegister godoc
// @Summary Register new Account
// @Description Register new Account
// @Param data body model.RegistrationAPI true "Account data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /accounts/register [post]
// @Tags Account
func PostAccountRegister(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())

	api := new(model.RegistrationAPI)
	if err := lib.BodyParser(c, api); err != nil {
		return lib.ErrorBadRequest(c, err)
	}

	var user model.User
	lib.Merge(api, &user)
	user.Password = lib.Strptr(lib.PasswordEncrypt(lib.RevStr(api.Password), viper.GetString("SALT"), viper.GetString("AES")))
	user.IsAdmin = lib.Boolptr(false)

	if err := db.Create(&user).Error; err != nil {
		return lib.ErrorConflict(c, err)
	}

	return lib.OK(c, user)
}
