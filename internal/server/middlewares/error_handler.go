package middlewares

import (
	"errors"
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/config"
)

func ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			c.Status(fiberErr.Code)
		} else {
			c.Status(fiber.StatusInternalServerError)
		}

		if c.Response().StatusCode() >= 500 {
			log.Printf("Error: %+v\n%s\n", err, debug.Stack())

			if !config.Env.App.Debug {
				return errors.New("server error")
			}
		}
	}

	return err
}
