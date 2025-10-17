package main

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/wallet"
	"github.com/kadsin/sms-gateway/tests"
	"github.com/stretchr/testify/require"
)

type mockProvider struct {
	fail bool
}

func (m *mockProvider) Send(s messages.Sms) error {
	if m.fail {
		return errors.New("failed to send")
	}
	return nil
}

func TestMain(m *testing.M) {
	tests.TestMain(m)
}

func Test_Queue_SendSms(t *testing.T) {
	m1 := tests.CreateSmsMessage(uuid.MustParse("ad760a00-3a02-4678-94cb-3546891dc083"))
	sendSms(context.Background(), &mockProvider{fail: true}, m1)

	var sms1 analytics_models.SmsMessage
	container.Analytics().Model(&analytics_models.SmsMessage{}).Find(&sms1, "id", m1.Id)
	require.Equal(t, sms1.Status, analytics_models.SMS_FAILED)

	m2 := tests.CreateSmsMessage(uuid.MustParse("ad760a00-3a02-4678-94cb-3546891dc083"))
	sendSms(context.Background(), &mockProvider{fail: false}, m2)

	var sms2 analytics_models.SmsMessage
	container.Analytics().Model(&analytics_models.SmsMessage{}).Find(&sms2, "id", m2.Id)
	require.Equal(t, sms2.Status, analytics_models.SMS_SENT)
}

func Test_Queue_RefundOnFailure(t *testing.T) {
	user := tests.CreateUser()
	oldBalance, _ := wallet.Get(context.Background(), user.ID)

	m := tests.CreateSmsMessage(user.ID)
	sendSms(context.Background(), &mockProvider{fail: true}, m)

	newBalance, _ := wallet.Get(context.Background(), user.ID)
	require.Equal(t, oldBalance+m.Price, newBalance)
}
