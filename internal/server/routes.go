package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/swagger"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/server/handlers"
)

func SetupRoutes(app *fiber.App) {
	setupSwagger(app.Group("/docs"))

	api := app.Group("/api")

	api.Post("/user/balance", handlers.ChangeUserBalance)

	api.Post("/sms", handlers.SendSms)

	api.Get("/reports", handlers.Reports)
}

func setupSwagger(router fiber.Router) {
	router.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			config.Env.Doc.Auth.Username: config.Env.Doc.Auth.Password,
		},
	}))

	router.Static("/swagger.yml", "./docs/swagger.yml")

	router.Get("/*", swagger.New(swagger.Config{
		Title: config.Env.App.Name + " - API Doc",
		URL:   "/docs/swagger.yml",
	}))
}
