package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	type User struct {
		ID      uuid.UUID `gorm:"primarykey; type: uuid; default: gen_random_uuid();"`
		Balance float32   `gorm:"not null; comment: Toman"`
		Email   string    `gorm:"type: varchar(320); unique; not null;"`

		gorm.Model
	}

	return container.DB().Migrator().CreateTable(&User{})
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	return container.DB().Migrator().DropTable("users")
}
