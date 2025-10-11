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
