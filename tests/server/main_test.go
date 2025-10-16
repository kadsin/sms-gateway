package server_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kadsin/sms-gateway/bootstrap"
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
