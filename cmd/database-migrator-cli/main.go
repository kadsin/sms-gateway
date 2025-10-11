package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	// It should load to execute all migration's init functions
	"github.com/kadsin/sms-gateway/config"
	"github.com/kadsin/sms-gateway/database"
	_ "github.com/kadsin/sms-gateway/database/migrations"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm/logger"
)

func main() {
	command, arguments := os.Args[1], os.Args[2:]

	db, err := goose.OpenDBWithDriver(config.Env.DB.Connection, database.DSN)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	if _, err := database.Connect(); err != nil {
		log.Printf("Error on connecting to database: %+v\n", err)
	}
	database.Instance().Logger.LogMode(logger.Silent)

	_, thisFile, _, _ := runtime.Caller(0)
	migrationsPath := fmt.Sprintf("%v/../../database/migrations", filepath.Dir(thisFile))

	if err := goose.RunContext(context.Background(), command, db, migrationsPath, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
