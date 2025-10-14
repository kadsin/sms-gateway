package database

import (
	"fmt"

	"github.com/kadsin/sms-gateway/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
	config.Env.DB.Host, config.Env.DB.Username, config.Env.DB.Password, config.Env.DB.Name, config.Env.DB.Port,
)

func New() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(DSN), &gorm.Config{})
}
