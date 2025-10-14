package container

import "log"

func Init() {
	initDB()

	initAnalytics()

	initKafkaProducer()
}

func Close() {
	db, err := DB().DB()
	if err == nil {
		db.Close()
	}

	analytics, err := Analytics().DB()
	if err == nil {
		analytics.Close()
	}

	KafkaProducer().Close()
}

func fatalOnNilInstance(msg string) {
	log.Fatal(msg + " Call container.Init() first.")
}
