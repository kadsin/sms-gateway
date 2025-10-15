package analytics_refresher

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"sync"

	// It should load to execute all migration's init functions
	_ "github.com/kadsin/sms-gateway/analytics/migrations"
	"github.com/kadsin/sms-gateway/internal/container"

	"github.com/kadsin/sms-gateway/analytics"
	"github.com/kadsin/sms-gateway/config"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm/logger"
)

var refreshDatabaseMux sync.Mutex

func RefreshDatabase() {
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

	MigrateDatabase()
}

func MigrateDatabase() {
	goose.SetLogger(goose.NopLogger())

	db, err := goose.OpenDBWithDriver(config.Env.Analytics.Connection, analytics.DSN)
	if err != nil {
		log.Fatalf("goose: failed to open analytics DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close analytics DB: %v\n", err)
		}
	}()

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsPath := fmt.Sprintf("%v/../../analytics/migrations", filepath.Dir(thisFile))

	if err := goose.RunContext(context.Background(), "up", db, migrationsPath); err != nil {
		log.Printf("Error on goose up: %v", err)
	}
}
