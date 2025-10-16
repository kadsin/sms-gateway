package server_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
)

func Test_Reports_Successful(t *testing.T) {
	user1 := tests.CreateUser(10000.5)
	generateBatchSms(user1, 10, 50, 150)

	user2 := tests.CreateUser()
	generateBatchSms(user2, 50, 10, 10)

	jsonBody := fmt.Sprintf(`{"client_id": "%s"}`, user1.ID)

	req := httptest.NewRequest(fiber.MethodGet, "/api/reports", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.StatusOK)

	type Payload struct {
		Data struct {
			Balance      float32 `json:"balance"`
			PendingCount int     `json:"pending_count"`
			FailureCount int     `json:"failure_count"`
			SentCount    int     `json:"sent_count"`
		} `json:"data"`
	}

	responseBody := unmarshalResponseBody[Payload](resp)

	require.Equal(t, float32(10000.5), responseBody.Data.Balance)
	require.Equal(t, 10, responseBody.Data.PendingCount)
	require.Equal(t, 50, responseBody.Data.FailureCount)
	require.Equal(t, 150, responseBody.Data.SentCount)
}

func generateBatchSms(user models.User, pendingCount, failureCount, sentCount int) {
	var pendings []*analytics_models.SmsMessage
	var failures []*analytics_models.SmsMessage
	var sents []*analytics_models.SmsMessage

	for range pendingCount {
		s := tests.CreateSmsMessage(user.ID)
		pendings = append(pendings, messageToModel(s, analytics_models.SMS_PENDING))
	}

	for range failureCount {
		s := tests.CreateSmsMessage(user.ID)
		pendings = append(pendings, messageToModel(s, analytics_models.SMS_PENDING))

		s.UpdatedAt = time.Now().Add(time.Hour)
		failures = append(failures, messageToModel(s, analytics_models.SMS_FAILED))
	}

	for range sentCount {
		s := tests.CreateSmsMessage(user.ID)
		pendings = append(pendings, messageToModel(s, analytics_models.SMS_PENDING))

		s.UpdatedAt = time.Now().Add(time.Hour)
		sents = append(sents, messageToModel(s, analytics_models.SMS_SENT))
	}

	container.Analytics().CreateInBatches(pendings, len(pendings))
	container.Analytics().CreateInBatches(failures, len(failures))
	container.Analytics().CreateInBatches(sents, len(sents))
}

func messageToModel(m messages.Sms, status analytics_models.SmsStatus) *analytics_models.SmsMessage {
	return &analytics_models.SmsMessage{
		ID:             m.Id,
		SenderClientID: m.SenderClientId,
		ReceiverPhone:  m.ReceiverPhone,
		Content:        m.Content,
		Price:          m.Price,
		IsExpress:      m.IsExpress,
		Status:         status,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
