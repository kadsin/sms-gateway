package sms

import (
	"github.com/kadsin/sms-gateway/internal/dtos/messages"
)

type Provider interface {
	Send(m messages.Sms) error
}
