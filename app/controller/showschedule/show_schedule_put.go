package showschedule

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PutShowSchedule godoc
// @Summary Update Show Schedule by id
// @Description Update Show Schedule by id
// @Param id path string true "Show Schedule ID"
// @Param data body model.ShowScheduleAPI true "Show Schedule data"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.ShowSchedule "Show Schedule data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /show-schedules/{id} [put]
// @Tags ShowSchedule
func PutShowSchedule(c *fiber.Ctx) error {
	if !lib.GetXIsAdmin(c) {
		return lib.ErrorUnauthorized(c)
	}

	api := new(model.ShowScheduleAPI)
	if err := lib.BodyParser(c, api); nil != err {
		return lib.ErrorBadRequest(c, err)
	}

	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.ShowSchedule
	result := db.Model(&data).
		Where(db.Where(model.ShowSchedule{
			Base: model.Base{
				ID: &id,
			},
		})).
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	lib.Merge(api, &data)
	data.ModifierID = lib.GetXUserID(c)

	if err := db.Model(&data).Updates(&data).Error; nil != err {
		return lib.ErrorConflict(c, err.Error())
	}

	return lib.OK(c, data)
}
