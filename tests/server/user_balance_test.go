package server_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/database"
	"github.com/kadsin/sms-gateway/database/models"
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

	tx := database.Instance().First(&models.User{}, "email = ?", "bad")

	require.ErrorIs(t, tx.Error, gorm.ErrRecordNotFound)
}

func Test_UserBalance_Success(t *testing.T) {
	jsonBody := `{
		"email": "test@example.com",
		"balance": 100000
	}`

	req := httptest.NewRequest(fiber.MethodPost, "/api/user/balance", bytes.NewReader([]byte(jsonBody)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	require.Equal(t, resp.StatusCode, fiber.StatusOK)

	var user models.User
	tx := database.Instance().First(&user, "email = ?", "test@example.com")

	require.NotErrorIs(t, tx.Error, gorm.ErrRecordNotFound)
	require.NotEqual(t, user.ID.String(), "00000000-0000-0000-0000-000000000000")
	require.Equal(t, user.Balance, float32(100000))
}
