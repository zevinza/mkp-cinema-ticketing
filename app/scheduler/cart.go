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
	redis.Set(context.Background(), "cart_scheduler", "off", 0)

	for range ticker.C {
		now := time.Now().Unix()

		carts := []model.Cart{}
		db.Where(`expires_in < ?`, now).Find(&carts)

		var ids []uuid.UUID

		if len(carts) > 0 {
			for _, cart := range carts {
				if cart.ID != nil {
					ids = append(ids, *cart.ID)
				}
			}
		} else {
			redis.Set(context.Background(), "cart_scheduler", "off", 0)
		}

		db.Where(`id IN ?`, ids).Unscoped().Delete(&model.Cart{})
	}
}
