package scheduler

import (
	"api/app/model"
	"api/app/services"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func CartScheduler() {
	ticker := time.NewTicker(time.Second * time.Duration(viper.GetInt64("CART_TICKER")))
	db := services.DB
	redis := services.REDIS
	if val, _ := redis.Get(context.Background(), "cart_scheduler").Result(); len(val) == 0 {
		redis.Set(context.Background(), "cart_scheduler", "off", 0)
	}

	for range ticker.C {
		if val, _ := redis.Get(context.Background(), "cart_scheduler").Result(); val == "on" {
			now := time.Now().Unix()

			carts := []model.Cart{}
			db.Where(`expires_in < ?`, now).Find(&carts)

			var ids []uuid.UUID

			if len(carts) > 0 {
				for _, cart := range carts {
					if cart.ID != nil {
						ids = append(ids, *cart.ID)
					}
					seatIDs := cart.GetSeat()
					db.Model(&model.Seat{}).Where(`id IN ?`, seatIDs).UpdateColumn("is_available", "true")
				}
				db.Where(`id IN ?`, ids).Unscoped().Delete(&model.Cart{})
			}

			var count int64
			db.Model(&model.Cart{}).Count(&count)
			if count == 0 {
				redis.Set(context.Background(), "cart_scheduler", "off", 0)
			}
		}
	}
}
