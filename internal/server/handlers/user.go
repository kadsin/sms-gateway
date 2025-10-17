package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/server/requests"
	"github.com/kadsin/sms-gateway/internal/wallet"
)

func ChangeUserBalance(c *fiber.Ctx) error {
	data, err := requests.Prepare[requests.UserBalanceRequest](c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	tx := container.DB().Begin()

	var user models.User
	if container.DB().Where("email", data.Email).FirstOrCreate(&user).Error != nil {
		return tx.Error
	}

	if _, err := wallet.Change(c.Context(), user.ID, data.Balance); err != nil {
		return err
	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id": user.ID,
	})
}
