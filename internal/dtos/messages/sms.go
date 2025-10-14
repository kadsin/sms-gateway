package messages

import "time"

type SmsStatus string

const (
	SMS_PENDING SmsStatus = "pending"
	SMS_SENT    SmsStatus = "sent"
	SMS_FAILED  SmsStatus = "failed"
)

type Sms struct {
	Id               string    `json:"id"`
	SenderClientId   string    `json:"sender_client_id"`
	ReceiverPhone    string    `json:"receiver_phone"`
	Content          string    `json:"content"`
	Price            float32   `json:"price"`
	IsExpress        bool      `json:"is_express"`
	Status           SmsStatus `json:"status"`
	TransferDuration float32   `json:"transfer_duration"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
