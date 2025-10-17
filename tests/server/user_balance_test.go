package server_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/wallet"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func Test_UserBalance_BadData(t *testing.T) {
	jsonBody := `{
		"email": "bad",
		"balance": "bad"
	}`

	req := httptest.NewRequest(fiber.MethodPost, "/api/user/balance", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.ErrUnprocessableEntity.Code)

	tx := container.DB().First(&models.User{}, "email", "bad")

	require.ErrorIs(t, tx.Error, gorm.ErrRecordNotFound)
}

func Test_UserBalance_NewUser(t *testing.T) {
	jsonBody := `{
		"email": "test@example.com",
		"balance": 100000
	}`

	req := httptest.NewRequest(fiber.MethodPost, "/api/user/balance", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.StatusOK)

	type Payload struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	responseBody := unmarshalResponseBody[Payload](resp)
	require.Regexp(t, "^.{8}-.{4}-.{4}-.{4}-.{12}$", responseBody.Data.ID)

	var user models.User
	tx := container.DB().First(&user, "email", "test@example.com")

	require.NotErrorIs(t, tx.Error, gorm.ErrRecordNotFound)
	require.NotEqual(t, user.ID.String(), "00000000-0000-0000-0000-000000000000")

	balance, _ := wallet.Get(context.Background(), user.ID)
	require.Equal(t, balance, float32(100000))
}

func Test_UserBalance_UpdateOldUser(t *testing.T) {
	user := tests.CreateUser(1500)
	wallet.Get(context.Background(), user.ID)

	jsonBody := fmt.Sprintf(`{
		"email": "%s",
		"balance": 200000
	}`, user.Email)

	req := httptest.NewRequest(fiber.MethodPost, "/api/user/balance", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	app.Test(req)

	var users []models.User
	container.DB().Model(&models.User{}).Where("email", user.Email).Scan(&users)

	require.Len(t, users, 1)

	balance, _ := wallet.Get(context.Background(), user.ID)
	require.Equal(t, float32(201500), balance)
}
