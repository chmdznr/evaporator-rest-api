package routes

import (
	"be-evaporator/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// Ping
	app.Get("/api/ping", controllers.GetPing)

	// Evaporator
	app.Post("/api/eva-data", controllers.CreateEvaData)
	app.Post("/api/eva-cv", controllers.CreateEvaCv)
}
