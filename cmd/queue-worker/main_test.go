package main

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
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
	m1 := message(uuid.MustParse("ad760a00-3a02-4678-94cb-3546891dc083"))
	sendSms(&mockProvider{fail: true}, m1)

	var sms1 analytics_models.SmsMessage
	container.Analytics().Model(&analytics_models.SmsMessage{}).Find(&sms1, "id", m1.Id)
	require.Equal(t, sms1.Status, analytics_models.SMS_FAILED)

	m2 := message(uuid.MustParse("ad760a00-3a02-4678-94cb-3546891dc083"))
	sendSms(&mockProvider{fail: false}, m2)

	var sms2 analytics_models.SmsMessage
	container.Analytics().Model(&analytics_models.SmsMessage{}).Find(&sms2, "id", m2.Id)
	require.Equal(t, sms2.Status, analytics_models.SMS_SENT)
}

func Test_Queue_IncreaseBalanceOnFailure(t *testing.T) {
	user := tests.CreateUser()

	m := message(user.ID)
	sendSms(&mockProvider{fail: true}, m)

	oldBalance := user.Balance
	container.DB().Find(&user, "id", user.ID)
	require.Equal(t, oldBalance+m.Price, user.Balance)
}

func message(sender uuid.UUID) messages.Sms {
	id, _ := uuid.NewV7()

	return messages.Sms{
		Id:             id,
		ReceiverPhone:  "+989156251013",
		SenderClientId: sender,
		Content:        "Temporary code: 1234",
		IsExpress:      false,
		Price:          105.6,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
