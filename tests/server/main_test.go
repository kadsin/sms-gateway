package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/bootstrap"
	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/tests"
)

var app *fiber.App

func TestMain(m *testing.M) {
	app = bootstrap.SetupFiberApp()

	tests.TestMain(m)
}

func unmarshalResponseBody[Payload any](resp *http.Response) (p Payload) {
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	json.Unmarshal(bytes, &p)
	return
}

func createUser(balance ...float32) models.User {
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
