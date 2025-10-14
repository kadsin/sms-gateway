package container

import (
	"log"

	"github.com/kadsin/sms-gateway/database"

	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func initDB() {
	var err error

	dbInstance, err = database.New()

	if err != nil {
		log.Fatalf("Error on initiating new db instance: %v", err)
	}
}

func DB() *gorm.DB {
	if dbInstance == nil {
		fatalOnNilInstance("Database is not connected.")
	}

	return dbInstance
}
