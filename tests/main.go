package tests

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	// It should load to execute all migration's init functions
	_ "github.com/kadsin/sms-gateway/analytics/migrations"
	_ "github.com/kadsin/sms-gateway/database/migrations"

	"github.com/kadsin/sms-gateway/analytics"
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/database"
	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/kadsin/sms-gateway/tests/mocks"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm/logger"
)

var refreshDatabaseMux sync.Mutex

var allGooseMigrations []*goose.Migration

func TestMain(m *testing.M) {
	container.Init()
	defer container.Close()

	mockEssentials()

	var err error
	if allGooseMigrations, err = goose.CollectMigrations(".", 0, math.MaxInt64); err != nil {
		log.Fatalf("error on collect goose migrations: %+v", err)
	}
	RefreshDatabase()
	RefreshAnalyticsDatabase()

	container.DB().Begin()
	defer container.DB().Rollback()

	container.Analytics().Begin()
	defer container.Analytics().Rollback()

	code := m.Run()

	os.Exit(code)
}

func mockEssentials() {
	container.MockKafkaProducer(&mocks.KafkaProducerMock{})

	container.MockKafkaConsumer(&mocks.KafkaConsumerMock{Topic: config.Env.Kafka.Topics.Regular})
}

func RefreshDatabase() {
	migrationsPath := collectAllGooseMigrations("database")

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

	migrate(config.Env.DB.Connection, database.DSN, migrationsPath)
}

func RefreshAnalyticsDatabase() {
	migrationsPath := collectAllGooseMigrations("analytics")

	refreshDatabaseMux.Lock()
	defer refreshDatabaseMux.Unlock()

	container.Analytics().Logger = logger.Discard

	if err := container.Analytics().Exec("DROP DATABASE IF EXISTS " + config.Env.Analytics.Name).Error; err != nil {
		log.Printf("failed to drop database: %v", err)
		return
	}

	if err := container.Analytics().Exec("CREATE DATABASE " + config.Env.Analytics.Name).Error; err != nil {
		log.Printf("failed to create database: %v", err)
		return
	}

	migrate(config.Env.Analytics.Connection, analytics.DSN, migrationsPath)
}

func collectAllGooseMigrations(scope string) (migrationsPath string) {
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsPath = fmt.Sprintf("%s/../%s/migrations", filepath.Dir(thisFile), scope)

	var migrations []*goose.Migration

	for _, m := range allGooseMigrations {
		if strings.Contains(m.Source, fmt.Sprintf("/%s/migrations/", scope)) {
			migrations = append(migrations, m)
		}
	}

	goose.ResetGlobalMigrations()
	goose.SetGlobalMigrations(migrations...)
	return
}

func migrate(driver, dsn, migrationsPath string) {
	goose.SetLogger(goose.NopLogger())

	db, err := goose.OpenDBWithDriver(driver, dsn)
	if err != nil {
		log.Fatalf("goose: failed to open %s: %v\n", driver, err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close %s: %v\n", driver, err)
		}
	}()

	if err := goose.RunContext(context.Background(), "up", db, migrationsPath); err != nil {
		log.Printf("Error on goose up for %s: %v", driver, err)
	}
}
