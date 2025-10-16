package userbalance

import (
	"context"
	"testing"

	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/kadsin/sms-gateway/tests/mocks"
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

func Test_UserBalance_Change_PublishMessage(t *testing.T) {
	container.KafkaProducer().(*mocks.KafkaProducerMock).Flush()

	c := container.KafkaConsumer("")
	c.(*mocks.KafkaConsumerMock).Topic = config.Env.Kafka.Topics.Balance

	ctx := context.Background()

	u := tests.CreateUser(5000)
	Change(ctx, u.ID, -100.57)

	m, _ := c.FetchMessage(ctx)
	b, _ := dtos.Unmarshal[messages.UserBalance](m.Value)

	require.Equal(t, u.ID.String(), b.ClientId.String())
	require.Equal(t, float32(-100.57), b.Amount)
}
