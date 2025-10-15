package messages

import (
	"time"

	"github.com/google/uuid"
)

type Sms struct {
	Id             uuid.UUID `json:"id"`
	SenderClientId uuid.UUID `json:"sender_client_id"`
	ReceiverPhone  string    `json:"receiver_phone"`
	Content        string    `json:"content"`
	Price          float32   `json:"price"`
	IsExpress      bool      `json:"is_express"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
