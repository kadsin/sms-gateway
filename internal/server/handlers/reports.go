package handlers

import (
	"github.com/gofiber/fiber/v2"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
)

func Reports(c *fiber.Ctx) error {
	var user models.User
	err := container.DB().Where("id", getClientId(c)).First(&user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "user not found")
	}

	stats, err := analytics_models.SmsStatsByClient(user.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"balance":       user.Balance,
		"pending_count": stats.Pending,
		"failure_count": stats.Failed,
		"sent_count":    stats.Sent,
	})
}
