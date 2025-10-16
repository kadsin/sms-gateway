package main

import (
	"context"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests.TestMain(m)
}

func Test_BalanceWorker_UpdateUserBalance(t *testing.T) {
	u := tests.CreateUser(10000.5)

	updateBalanceInPostgres(context.Background(), messages.UserBalance{
		ClientId: u.ID,
		Amount:   -1000.4,
	})

	container.DB().Where("id", u.ID).First(&u)
	require.Equal(t, float32(9000.1), u.Balance)
}
