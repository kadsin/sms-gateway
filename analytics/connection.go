package analytics

import (
	"fmt"
	"log"

	"github.com/kadsin/sms-gateway/config"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var db *gorm.DB

var DSN = fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s",
	config.Env.Analytics.Username, config.Env.Analytics.Password, config.Env.Analytics.Host, config.Env.Analytics.Port, config.Env.Analytics.Name,
)

func Connect() (*gorm.DB, error) {
	var err error

	db, err = gorm.Open(clickhouse.Open(DSN), &gorm.Config{})

	return db, err
}

func Instance() *gorm.DB {
	if db == nil {
		log.Fatal("Database not connected. Call database.Connect() first.")
	}
	return db
}
