package routes

import (
	"api/app/controller"
	"api/app/controller/account"
	"api/app/controller/cart"
	"api/app/controller/showschedule"
	"api/app/controller/transaction"
	"api/app/controller/user"
	"api/app/lib"
	"api/app/middleware"
	"api/app/scheduler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
)

// Handle all request to route to controller
func Handle(app *fiber.App) {
	app.Use(cors.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			lib.PrintStackTrace(e)
		},
	}))

	go scheduler.CartScheduler()

	api := app.Group(viper.GetString("ENDPOINT"))

	api.Static("/swagger", "docs/swagger.json")
	api.Get("/", controller.GetAPIIndex)
	api.Get("/info.json", controller.GetAPIInfo)
	api.Post("/logs", controller.PostLogs)

	// Account
	accountAPI := api.Group("/accounts")
	accountAPI.Use(middleware.TokenValidator())
	accountAPI.Post("/register", account.PostAccountRegister)
	accountAPI.Post("/login", account.PostAccountLogin)
	accountAPI.Post("/logout", account.PostAccountLogout)
	accountAPI.Post("/refresh", account.PostAccountRefresh)

	// Cart
	cartAPI := api.Group("/carts")
	cartAPI.Use(middleware.Oauth2Authentication)
	cartAPI.Post("/", cart.PostCart)
	cartAPI.Get("/", cart.GetCart)
	cartAPI.Put("/:id", cart.PutCart)
	cartAPI.Get("/:id", cart.GetCartID)
	cartAPI.Delete("/:id", cart.DeleteCart)

	// Show Schedule
	showscheduleAPI := api.Group("/show-schedules")
	showscheduleAPI.Use(middleware.TokenValidator())
	showscheduleAPI.Post("/", showschedule.PostShowSchedule)
	showscheduleAPI.Get("/", showschedule.GetShowSchedule)
	showscheduleAPI.Put("/:id", showschedule.PutShowSchedule)
	showscheduleAPI.Get("/:id", showschedule.GetShowScheduleID)
	showscheduleAPI.Delete("/:id", showschedule.DeleteShowSchedule)

	// transaction
	transactionAPI := api.Group("/transactions")
	transactionAPI.Use(middleware.Oauth2Authentication)
	transactionAPI.Post("/", transaction.PostTransaction)
	transactionAPI.Get("/", transaction.GetTransaction)
	transactionAPI.Get("/:id", transaction.GetTransactionID)
	transactionAPI.Post("/:id/payment", transaction.PostTransactionPayment)
	transactionAPI.Post("/:id/cancel", transaction.PostTransactionCancel)

	// User
	userAPI := api.Group("/users")
	userAPI.Use(middleware.Oauth2Authentication)
	userAPI.Get("/", user.GetUser)
	userAPI.Put("/:id", user.PutUser)
	userAPI.Get("/:id", user.GetUserID)
	userAPI.Post("/", user.PostUser)

}
