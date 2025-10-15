package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	analytics_models "github.com/kadsin/sms-gateway/analytics/models"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSmsMessagesTable, downCreateSmsMessagesTable)
}

func upCreateSmsMessagesTable(ctx context.Context, tx *sql.Tx) error {
	type SmsMessage struct {
		ID             uuid.UUID `gorm:"primaryKey"`
		SenderClientID uuid.UUID
		ReceiverPhone  string
		Content        string
		Price          float32
		IsExpress      bool
		Status         analytics_models.SmsStatus
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	return container.Analytics().Migrator().CreateTable(&SmsMessage{})
}

func downCreateSmsMessagesTable(ctx context.Context, tx *sql.Tx) error {
	return container.Analytics().Migrator().DropTable("sms_messages")
}
