package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSmsMessagesTable, downCreateSmsMessagesTable)
}

func upCreateSmsMessagesTable(ctx context.Context, tx *sql.Tx) error {
	ddl := `
	CREATE TABLE IF NOT EXISTS sms_messages
	(
		id UUID,
		sender_client_id UUID,
		receiver_phone String,
		content String,
		price Float32,
		is_express UInt8,
		status String,
		created_at DateTime,
		updated_at DateTime
	) ENGINE = MergeTree()
	PARTITION BY toYYYYMMDD(created_at)
	ORDER BY (sender_client_id, id)
	`

	_, err := tx.ExecContext(ctx, ddl)
	return err
}

func downCreateSmsMessagesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS sms_messages`)
	return err
}
