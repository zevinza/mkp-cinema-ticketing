package showschedule

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
)

// PostShowSchedule godoc
// @Summary Create new Show Schedule
// @Description Create new Show Schedule
// @Param data body model.ShowScheduleAPI true "Show Schedule data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.ShowSchedule "Show Schedule data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /show-schedules [post]
// @Tags ShowSchedule
func PostShowSchedule(c *fiber.Ctx) error {
	if !lib.GetXIsAdmin(c) {
		return lib.ErrorUnauthorized(c)
	}

	api := new(model.ShowScheduleAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	db := services.DB.WithContext(c.UserContext())

	var data model.ShowSchedule
	lib.Merge(api, &data)
	data.CreatorID = lib.GetXUserID(c)

	if err := db.Create(&data).Error; nil != err {
		return lib.ErrorConflict(c, err.Error())
	}

	return lib.Created(c, data)
}
