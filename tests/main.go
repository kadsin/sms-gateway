package tests

import (
	"os"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
)

func TestMain(m *testing.M) {
	container.Init()
	defer container.Close()

	RefreshDatabase()

	container.DB().Begin()
	code := m.Run()
	container.DB().Rollback()

	os.Exit(code)
}
