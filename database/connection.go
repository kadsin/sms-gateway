package database

import (
	"fmt"
	"log"

	"github.com/kadsin/sms-gateway/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
	config.Env.DB.Host, config.Env.DB.Username, config.Env.DB.Password, config.Env.DB.Name, config.Env.DB.Port,
)

func Connect() (*gorm.DB, error) {
	var err error

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})

	return db, err
}

func Instance() *gorm.DB {
	if db == nil {
		log.Fatal("Database not connected. Call database.Connect() first.")
	}
	return db
}
