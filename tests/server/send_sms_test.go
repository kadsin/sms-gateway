package server_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/kadsin/sms-gateway/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SendSms_BadData(t *testing.T) {
	user := tests.CreateUser()

	smsContent := strings.Repeat("a", 161)

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "091",
			"content": "%s",
			"is_express": 10
		}`, user.Email, smsContent)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.ErrUnprocessableEntity.Code)

	c := container.KafkaConsumer("")

	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Regular
	_, err := c.FetchMessage(context.Background())
	require.NotNil(t, err)

	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Express
	_, err = c.FetchMessage(context.Background())
	require.NotNil(t, err)
}

func Test_SendSms_CalculatePrice(t *testing.T) {
	user := tests.CreateUser()

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "+989123456789",
			"content": "abcd",
			"is_express": false
		}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.StatusOK)

	c := container.KafkaConsumer("")
	defer container.KafkaProducer().Close()

	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Regular
	m, _ := c.FetchMessage(context.Background())

	sms, _ := dtos.Unmarshal[messages.Sms](m.Value)
	require.Equal(t, float32(40), sms.Price)

	oldBalance := user.Balance
	container.DB().Find(&user, "id", user.ID)
	require.Equal(t, oldBalance-40, user.Balance)
}

func Test_SendSms_ValidateBalance(t *testing.T) {
	user := tests.CreateUser(10)

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "+989123456789",
			"content": "abcdefg",
			"is_express": false
		}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.ErrPaymentRequired.Code)

	c := container.KafkaConsumer("")
	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Regular

	_, err := c.FetchMessage(context.Background())
	require.NotNil(t, err)
}

func Test_SendSms_Express(t *testing.T) {
	user := tests.CreateUser()

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "+989123456789",
			"content": "abcdefg",
			"is_express": true
		}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	app.Test(req)

	c := container.KafkaConsumer("")
	defer container.KafkaProducer().Close()
	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Express

	_, err := c.FetchMessage(context.Background())
	require.Nil(t, err)
}

func Test_SendSms_SuccessfulResponse(t *testing.T) {
	user := tests.CreateUser()

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "+989123456789",
			"content": "aqfdvsvsdvs",
			"is_express": false
		}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	c := container.KafkaConsumer("")
	defer container.KafkaProducer().Close()
	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Regular

	m, _ := c.FetchMessage(context.Background())
	sms, _ := dtos.Unmarshal[messages.Sms](m.Value)

	type Payload struct {
		Data struct {
			MessageId string  `json:"message_id"`
			Price     float32 `json:"price"`
		} `json:"data"`
	}

	responseBody := unmarshalResponseBody[Payload](resp)

	require.Equal(t, sms.Id.String(), responseBody.Data.MessageId)
	require.Equal(t, sms.Price, responseBody.Data.Price)
}

func Test_SendSms_SuccessfulStoreInClickHouse(t *testing.T) {
	user := tests.CreateUser()

	jsonBody := fmt.Sprintf(`{
			"client_email": "%s",
			"receiver_phone": "+989123456789",
			"content": "aqfdvsvsdvs",
			"is_express": false
		}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/sms", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	type Payload struct {
		Data struct {
			MessageId string  `json:"message_id"`
			Price     float32 `json:"price"`
		} `json:"data"`
	}

	responseBody := unmarshalResponseBody[Payload](resp)

	var sms analytics_models.SmsMessage
	container.Analytics().Model(&analytics_models.SmsMessage{}).First(&sms, "id", responseBody.Data.MessageId)
	assert.Equal(t, sms.Content, "aqfdvsvsdvs")
	assert.Equal(t, sms.Status, analytics_models.SMS_PENDING)
}
