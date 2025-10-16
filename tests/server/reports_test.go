package server_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
)

func Test_Reports_Successful(t *testing.T) {
	user1 := tests.CreateUser()
	generateBatchSms(user1, 10, 50, 150)

	user2 := tests.CreateUser()
	generateBatchSms(user2, 50, 10, 10)

	jsonBody := fmt.Sprintf(`{"client_id": "%s"}`, user1.ID)

	req := httptest.NewRequest(fiber.MethodGet, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.StatusOK)

	type Payload struct {
		Data struct {
			Balance      float32 `json:"balance"`
			PendingCount float32 `json:"pending_count"`
			FailureCount float32 `json:"failure_count"`
			SentCount    float32 `json:"sent_count"`
		} `json:"data"`
	}

	responseBody := unmarshalResponseBody[Payload](resp)

	require.Equal(t, user1.Balance, responseBody.Data.Balance)
	require.Equal(t, 10, responseBody.Data.PendingCount)
	require.Equal(t, 50, responseBody.Data.FailureCount)
	require.Equal(t, 150, responseBody.Data.SentCount)
}

func generateBatchSms(user models.User, pendingCount, failureCount, sentCount int) {
	for range pendingCount {
		s := tests.CreateSmsMessage(user.ID)
		analytics_models.LogPending(&s)
	}

	for range failureCount {
		s := tests.CreateSmsMessage(user.ID)
		analytics_models.LogFailure(&s)
	}

	for range sentCount {
		s := tests.CreateSmsMessage(user.ID)
		analytics_models.LogSent(&s)
	}
}
