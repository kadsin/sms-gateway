package main

import (
	"github.com/kadsin/sms-gateway/bootstrap"
)

func main() {
	bootstrap.SetupDatabase()

	app := bootstrap.SetupFiberApp()

	app.Listen(":3000")
}
