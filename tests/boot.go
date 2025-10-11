package tests

import (
	"os"
	"testing"

	"github.com/kadsin/sms-gateway/database"
)

func TestMain(m *testing.M) {
	RefreshDatabase()

	database.Instance().Begin()
	code := m.Run()
	database.Instance().Rollback()

	os.Exit(code)
}
