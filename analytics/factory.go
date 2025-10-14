package analytics

import (
	"fmt"

	"github.com/kadsin/sms-gateway/config"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

var db *gorm.DB

var DSN = fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s",
	config.Env.Analytics.Username, config.Env.Analytics.Password, config.Env.Analytics.Host, config.Env.Analytics.Port, config.Env.Analytics.Name,
)

func New() (*gorm.DB, error) {
	return gorm.Open(clickhouse.Open(DSN), &gorm.Config{})
}
