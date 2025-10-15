package analytics_models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
)

type SmsStatus string

func (s SmsStatus) String() string {
	return string(s)
}

const (
	SMS_PENDING SmsStatus = "pending"
	SMS_SENT    SmsStatus = "sent"
	SMS_FAILED  SmsStatus = "failed"
)

type SmsMessage struct {
	ID             uuid.UUID `gorm:"primaryKey"`
	SenderClientID uuid.UUID
	ReceiverPhone  string
	Content        string
	Price          float32
	IsExpress      bool
	Status         SmsStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func LogPending(m *messages.Sms) {
	createSmsMessage(m, SMS_PENDING)
}

func LogFailure(m *messages.Sms) {
	m.UpdatedAt = time.Now()
	createSmsMessage(m, SMS_FAILED)
}

func LogSent(m *messages.Sms) {
	m.UpdatedAt = time.Now()
	createSmsMessage(m, SMS_SENT)
}

func createSmsMessage(m *messages.Sms, status SmsStatus) {
	tx := container.Analytics().Create(&SmsMessage{
		ID:             m.Id,
		SenderClientID: m.SenderClientId,
		ReceiverPhone:  m.ReceiverPhone,
		Content:        m.Content,
		Price:          m.Price,
		IsExpress:      m.IsExpress,
		Status:         status,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	})

	if tx.Error != nil {
		log.Printf("Error on create sms message: %+v", tx.Error)
	}
}
