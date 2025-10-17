package messages

import (
	"github.com/google/uuid"
)

type WalletChanged struct {
	ClientId uuid.UUID `json:"client_id"`
	Amount   float32   `json:"amount"`
}
