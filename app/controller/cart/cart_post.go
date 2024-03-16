package cart

import (
	"api/app/lib"
	"api/app/model"
	"api/app/services"
	"context"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// PostCart godoc
// @Summary Create new Cart
// @Description Create new Cart
// @Param data body model.CartPayload true "Cart data"
// @Accept  application/json
// @Produce application/json
// @Success 201 {object} model.Cart "Cart data"
// @Failure 400 {object} lib.Response
// @Failure 404 {object} lib.Response
// @Failure 409 {object} lib.Response
// @Failure 500 {object} lib.Response
// @Failure default {object} lib.Response
// @Security ApiKeyAuth
// @Router /carts [post]
// @Tags Cart
func PostCart(c *fiber.Ctx) error {
	api := new(model.CartPayload)
	if err := c.BodyParser(api); err != nil {
		return lib.ErrorBadRequest(c, err)
	}

	db := services.DB.WithContext(c.UserContext())
	rdb := services.REDIS
	ctx := context.Background()

	duration := time.Duration(time.Second * time.Duration(int64(viper.GetInt("CART_DURATION"))))

	userID := lib.GetXUserID(c)
	cartID := lib.GenUUID()

	var count int64
	db.Model(&model.Cart{}).Where(`user_id = ?`, userID).Count(&count)
	if count > 0 {
		return lib.ErrorBadRequest(c, "please proceed your previous order")
	}
	for _, id := range api.SeatIDs {
		key := "seat_" + id.String()
		if _, err := rdb.Get(ctx, key).Result(); err != nil {
			return lib.ErrorBadRequest(c, "your seat is taken, please choose another")
		}

		rdb.Set(ctx, key, cartID.String, duration)
	}
	go InsertCart(api, userID, cartID, db)
	go db.Model(&model.Seat{}).Where(`id IN ?`, api.SeatIDs).UpdateColumn("is_available", "false")

	rdb.Set(ctx, "cart_"+cartID.String(), "ok", duration)
	rdb.Set(ctx, "cart_scheduler", "on", 0)
	log.Println("-----seat scheduler = on-----")

	return lib.OK(c, cartID)
}

func InsertCart(api *model.CartPayload, userID, cartID *uuid.UUID, db *gorm.DB) {
	var ids string
	for _, id := range api.SeatIDs {
		ids = ids + ", " + id.String()
	}
	ids = strings.TrimRight(ids, ", ")

	show := model.ShowSchedule{}
	row := db.Where(`id = ?`, api.ShowScheduleID).Take(&show)
	if row.RowsAffected == 0 {
		return
	}
	qty := len(api.SeatIDs)
	total := lib.RevFloat64(show.Price) * float64(qty)

	db.Create(&model.Cart{
		Base: model.Base{
			ID: cartID,
		},
		CartAPI: model.CartAPI{
			UserID:         userID,
			ShowScheduleID: api.ShowScheduleID,
			SeatIDs:        lib.Strptr(ids),
			Price:          show.Price,
			Quantity:       &qty,
			TotalPrice:     &total,
		},
	})
}
