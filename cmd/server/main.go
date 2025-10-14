package main

import (
	"github.com/kadsin/sms-gateway/bootstrap"
	"github.com/kadsin/sms-gateway/internal/container"
)

func main() {
	container.Init()
	defer container.Close()

	app := bootstrap.SetupFiberApp()

	app.Listen(":3000")
}
