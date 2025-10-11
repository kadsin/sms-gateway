package bootstrap

import (
	"github.com/kadsin/sms-gateway/internal/server"
	"github.com/kadsin/sms-gateway/internal/server/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupFiberApp() *fiber.App {
	app := fiber.New()

	app.Use("/api",
		cors.New(),

		middlewares.ResponseWrapper,
		middlewares.ErrorHandler,

		recover.New(),
	)

	server.SetupRoutes(app)

	return app
}
