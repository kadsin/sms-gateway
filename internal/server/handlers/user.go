package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/internal/server/requests"
)

func ChangeUserBalance(c *fiber.Ctx) error {
	_, err := requests.Prepare[requests.UserBalanceRequest](c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	return nil
}
