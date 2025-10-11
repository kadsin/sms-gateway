package bootstrap

import (
	"log"

	"github.com/kadsin/sms-gateway/database"
)

func SetupDatabase() {
	_, err := database.Connect()
	if err != nil {
		log.Printf("Error on connect to database: %+v\n", err)
		panic(err)
	}
}
