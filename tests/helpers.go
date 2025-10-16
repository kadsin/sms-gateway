package tests

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
)

func CreateUser(balance ...float32) models.User {
	var b float32 = 100000.55
	if len(balance) > 0 {
		b = balance[0]
	}

	user := models.User{
		Email:   "test" + strconv.Itoa(time.Now().Nanosecond()) + "@example.com",
		Balance: b,
	}

	container.DB().Create(&user)
	return user
}

func CreateSmsMessage(sender uuid.UUID) messages.Sms {
	id, _ := uuid.NewV7()

	return messages.Sms{
		Id:             id,
		ReceiverPhone:  "+989156251013",
		SenderClientId: sender,
		Content:        "Temporary code: 1234",
		IsExpress:      false,
		Price:          105.6,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
