package container

import (
	"fmt"
	"math/rand/v2"

	"github.com/kadsin/sms-gateway/internal/dtos/messages"
	"github.com/kadsin/sms-gateway/internal/sms"
)

type fakeSmsProvider struct{}

func (fakeSmsProvider) Send(m messages.Sms) error {
	if sent := rand.IntN(2) == 1; !sent {
		return fmt.Errorf("message sending failed with id %v", m.Id)
	}

	return nil
}

func SmsProvider() sms.Provider {
	if mockSmsProviderInstance != nil {
		return mockSmsProviderInstance
	}

	return &fakeSmsProvider{}
}
