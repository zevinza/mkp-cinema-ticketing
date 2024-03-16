package showschedule

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetShowScheduleID godoc
// @Summary Get a Show Schedule by id
// @Description Get a Show Schedule by id
// @Param id path string true "Show Schedule ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} model.ShowSchedule "Show Schedule data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /show-schedules/{id} [get]
// @Tags ShowSchedule
func GetShowScheduleID(c *fiber.Ctx) error {
	db := services.DB.WithContext(c.UserContext())
	id, _ := uuid.Parse(c.Params("id"))

	var data model.ShowSchedule
	result := db.Model(&data).
		Where(db.Where(model.ShowSchedule{
			Base: model.Base{
				ID: &id,
			},
		})).
		Preload("Movie").
		Preload("CinemaLocation").
		Preload("Theater").
		Preload("Seats").
		Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	return lib.OK(c, data)
}
