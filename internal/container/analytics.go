package container

import (
	"log"

	"github.com/kadsin/sms-gateway/analytics"
	"gorm.io/gorm"
)

var analyticsInstance *gorm.DB

func initAnalytics() {
	var err error

	analyticsInstance, err = analytics.New()

	if err != nil {
		log.Fatalf("Error on initiating new analytics db instance: %v", err)
	}
}

func Analytics() *gorm.DB {
	if analyticsInstance == nil {
		fatalOnNilInstance("Analytics database is not connected.")
	}

	return analyticsInstance
}
