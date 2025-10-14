package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/server/requests"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const SmsPricePerChar float32 = 10 // Toman

func SendSms(c *fiber.Ctx) error {
	data, err := requests.Prepare[requests.SmsRequest](c)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	tx := container.DB().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	var user models.User
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("email", data.ClientEmail).
		First(&user).Error
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "user not found")
	}

	smsPrice := float32(len(data.Content)) * SmsPricePerChar
	if user.Balance < smsPrice {
		return fiber.NewError(
			fiber.StatusPaymentRequired,
			fmt.Sprintf("balance is not enough. price: %v, balance: %v", smsPrice, user.Balance),
		)
	}

	smsMessage, err := generateSmsMessage(data)
	if err != nil {
		return err
	}

	smsMessage.Price = smsPrice
	smsMessage.SenderClientId = user.ID.String()

	marshaledSms, err := dtos.Marshal(smsMessage)
	if err != nil {
		return err
	}

	kafkaTopic := config.Env.Kafka.Topics.Regular
	if *data.IsExpress {
		kafkaTopic = config.Env.Kafka.Topics.Express
	}

	err = container.KafkaProducer().SendMessage(c.Context(), kafka.Message{
		Topic: kafkaTopic,
		Value: marshaledSms,
	})
	if err != nil {
		return err
	}

	err = tx.Model(&models.User{}).
		Where("id", user.ID).
		UpdateColumn("balance", gorm.Expr("balance - ?", smsPrice)).Error
	if err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

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
		Id:            id.String(),
		ReceiverPhone: data.ReceiverPhone,
		Content:       data.Content,
		IsExpress:     *data.IsExpress,
		Status:        messages.SMS_PENDING,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	return sms, nil
}
