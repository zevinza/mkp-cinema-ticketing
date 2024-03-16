package main

import (
	"api/app/config"
	"api/app/lib"
	"api/app/routes"
	"api/app/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	lib.LoadEnvironment(config.Environment)
	services.InitDatabase()
	services.InitRedis()
}

// @title MKP Cinema Ticketing
// @version 1.0.0
// @description API Documentation
// @contact.name Armada Muhammad Siswanto
// @contact.email armadamuhammads@gmail.com
// @host localhost:8000
// @schemes http
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey TokenKey
// @in header
// @name x-Token
func main() {
	app := fiber.New(fiber.Config{
		Prefork: viper.GetString("PREFORK") == "true",
	})

	routes.Handle(app)
	log.Fatal(app.Listen(":" + viper.GetString("PORT")))
}
