package userbalance

import (
	"context"
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	tests.TestMain(m)
}

func Test_UserBalance_NotExistsOnRedis(t *testing.T) {
	ctx := context.Background()

	u := tests.CreateUser(5000)

	Get(ctx, u.ID)
	rb, _ := container.Redis().Get(ctx, CacheKey(u.ID)).Float32()

	require.Equal(t, 5000, rb)
}
