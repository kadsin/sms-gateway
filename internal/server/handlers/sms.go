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

const SmsPricePerChar float32 = 10 // Toman

func SendSms(c *fiber.Ctx) error {
	data, err := requests.Prepare[requests.SmsRequest](c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	userId, err := uuid.Parse(getClientId(c))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid user id")
	}

	userMux := wallet.UserLock(userId)
	userMux.Lock()
	defer userMux.Unlock()

	balance, err := wallet.Get(c.Context(), userId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Error on getting balance")
	}

	smsPrice := float32(len(data.Content)) * SmsPricePerChar

	if balance < smsPrice {
		return fiber.NewError(
			fiber.StatusPaymentRequired,
			fmt.Sprintf("balance is not enough. price: %v, balance: %v", smsPrice, balance),
		)
	}

	smsMessage, err := generateSmsMessage(data)
	if err != nil {
		return err
	}

	smsMessage.Price = smsPrice
	smsMessage.SenderClientId = userId

	marshaledSms, err := dtos.Marshal(smsMessage)
	if err != nil {
		return err
	}

	kafkaTopic := config.Env.Kafka.Topics.Regular
	if *data.IsExpress {
		kafkaTopic = config.Env.Kafka.Topics.Express
	}

	if _, err := wallet.Change(c.Context(), userId, -smsPrice); err != nil {
		return err
	}

	err = container.KafkaProducer().SendMessage(c.Context(), kafka.Message{
		Topic: kafkaTopic,
		Value: marshaledSms,
	})
	if err != nil {
		// Refund
		wallet.Change(c.Context(), userId, smsPrice)
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

	sms := &messages.Sms{
		Id:            id,
		ReceiverPhone: data.ReceiverPhone,
		Content:       data.Content,
		IsExpress:     *data.IsExpress,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return sms, nil
}
