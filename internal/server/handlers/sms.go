package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/server/requests"
	"github.com/kadsin/sms-gateway/internal/wallet"
	"github.com/segmentio/kafka-go"
)

const SMS_PRICE_PER_CHAR float32 = 10 // Toman

func SendSms(c *fiber.Ctx) error {
	data, err := requests.Prepare[requests.SmsRequest](c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	userId, err := uuid.Parse(getClientId(c))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid user id")
	}

	smsMessage, err := generateSmsMessage(data)
	if err != nil {
		return err
	}
	smsMessage.SenderClientId = userId

	marshaledSms, err := dtos.Marshal(smsMessage)
	if err != nil {
		return err
	}

	kafkaTopic := config.Env.Kafka.Topics.Regular
	if *data.IsExpress {
		kafkaTopic = config.Env.Kafka.Topics.Express
	}

	newBalance, err := wallet.Change(c.Context(), userId, -smsMessage.Price)
	if err != nil {
		if err == wallet.ErrInsufficientFunds {
			return fiber.NewError(
				fiber.StatusPaymentRequired,
				fmt.Sprintf("not enough balance (price %.2f, balance %.2f)", smsMessage.Price, newBalance),
			)
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = container.KafkaProducer().SendMessage(c.Context(), kafka.Message{
		Topic: kafkaTopic,
		Value: marshaledSms,
	})
	if err != nil {
		// Refund
		wallet.Change(c.Context(), userId, smsMessage.Price)
		return err
	}

	analytics_models.LogPending(smsMessage)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message_id": smsMessage.Id,
		"price":      smsMessage.Price,
	})
}

func generateSmsMessage(data requests.SmsRequest) (*messages.Sms, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &messages.Sms{
		Id:            id,
		ReceiverPhone: data.ReceiverPhone,
		Content:       data.Content,
		Price:         float32(len(data.Content)) * SMS_PRICE_PER_CHAR,
		IsExpress:     *data.IsExpress,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}
