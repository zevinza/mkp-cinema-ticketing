package account

import (
	"api/app/lib"
	"api/app/middleware"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// PostAccountLogin godoc
// @Summary Login into Account
// @Description Login into Account
// @Param data body model.LoginAPI true "Account data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.LoginResponse "Account data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security TokenKey
// @Router /accounts/login [post]
// @Tags Account
func PostAccountLogin(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())

	api := new(model.LoginAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	user := model.User{}
	if row := db.Where(`username = ? OR email = ?`, api.Username, api.Username).Take(&user); row.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	if ok := lib.PasswordCompare(lib.RevStr(user.Password), lib.RevStr(api.Password), viper.GetString("SALT"), viper.GetString("AES")); !ok {
		return lib.ErrorUnauthorized(c, "password didn't match")
	}

	res, err := middleware.GenerateAccessToken(&user)
	if nil != err {
		return lib.ErrorUnauthorized(c, err.Error())
	}

	return lib.OK(c, &model.LoginResponse{
		Token: &res,
		User:  &user,
	})

}
