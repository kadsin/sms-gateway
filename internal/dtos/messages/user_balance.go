package messages

import (
	"github.com/google/uuid"
)

type UserBalance struct {
	ClientId uuid.UUID `json:"client_id"`
	Amount   float32   `json:"amount"`
}
