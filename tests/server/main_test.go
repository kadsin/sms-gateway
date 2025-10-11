package server_test

import (
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
