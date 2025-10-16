package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func getClientId(c *fiber.Ctx) string {
	return c.Get("X-CLIENT-ID")
}
