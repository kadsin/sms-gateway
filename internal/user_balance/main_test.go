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

func Test_UserBalance_Get_NotExistsOnRedis(t *testing.T) {
	ctx := context.Background()

	u := tests.CreateUser(5000)

	Get(ctx, u.ID)
	rb, _ := container.Redis().Get(ctx, CacheKey(u.ID)).Float32()

	require.Equal(t, float32(5000), rb)
}

func Test_UserBalance_Increment_NotExistsOnRedis(t *testing.T) {
	ctx := context.Background()

	u := tests.CreateUser(5000)

	Change(ctx, u.ID, +1000.55)
	rb, _ := container.Redis().Get(ctx, CacheKey(u.ID)).Float32()

	require.Equal(t, float32(6000.55), rb)
}

func Test_UserBalance_Decrement(t *testing.T) {
	ctx := context.Background()

	u := tests.CreateUser(5000)

	Change(ctx, u.ID, -100.57)
	nb, _ := Get(ctx, u.ID)

	require.Equal(t, float32(4899.43), nb)
}
