package migrations

import (
	"context"
	"database/sql"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	type Example struct {
		gorm.Model

		Name string
	}

	return container.Analytics().Migrator().CreateTable(&Example{})
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	return container.Analytics().Migrator().DropTable("examples")
}
