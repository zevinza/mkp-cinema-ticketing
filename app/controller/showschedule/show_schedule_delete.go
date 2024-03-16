package showschedule

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DeleteShowSchedule godoc
// @Summary Delete Show Schedule by id
// @Description Delete Show Schedule by id
// @Param id path string true "Show Schedule ID"
// @Accept  application/json
// @Produce application/json
// @Success 200 {object} lib.Response
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /show-schedules/{id} [delete]
// @Tags ShowSchedule
func DeleteShowSchedule(c *fiber.Ctx) error {
	if !lib.GetXIsAdmin(c) {
		return lib.ErrorUnauthorized(c)
	}

	db := services.DB.WithContext(c.UserContext())
	id := c.Params("id")

	var data model.ShowSchedule
	result := db.Model(&data).Where("id = ?", id).Take(&data)
	if result.RowsAffected < 1 {
		return lib.ErrorNotFound(c)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := db.Delete(&data).Error; err != nil {
			return err
		}

		if err := db.Where(`show_schedule_id = ?`, id).Unscoped().Delete(&model.Seat{}).Error; err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		return lib.ErrorConflict(c, err)
	}

	return lib.OK(c)
}
