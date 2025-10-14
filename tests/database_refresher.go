package tests

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	// It should load to execute all migration's init functions
	_ "github.com/kadsin/sms-gateway/database/migrations"
	"github.com/kadsin/sms-gateway/internal/container"

	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/database"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm/logger"
)

var refreshDatabaseMux sync.Mutex

func RefreshDatabase() {
	refreshDatabaseMux.Lock()
	defer refreshDatabaseMux.Unlock()

	container.DB().Logger = logger.Discard

	tables, err := container.DB().Migrator().GetTables()
	if err != nil {
		log.Fatalf("Could not get tables list: %v", err)
	}

	for _, t := range tables {
		err = container.DB().Migrator().DropTable(t)
		if err != nil {
			log.Fatalf("Could not drop tables: %v", err)
		}
	}

	MigrateDatabase()
}

func MigrateDatabase() {
	goose.SetLogger(goose.NopLogger())

	db, err := goose.OpenDBWithDriver(config.Env.DB.Connection, database.DSN)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsPath := fmt.Sprintf("%v/../database/migrations", filepath.Dir(thisFile))

	if err := goose.RunContext(context.Background(), "up", db, migrationsPath); err != nil {
		log.Printf("Error on goose up: %v", err)
	}
}
