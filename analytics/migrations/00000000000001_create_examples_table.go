package migrations

import (
	"context"
	"database/sql"

	"github.com/kadsin/sms-gateway/analytics"
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

	return analytics.Instance().Migrator().CreateTable(&Example{})
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	return analytics.Instance().Migrator().DropTable("examples")
}
