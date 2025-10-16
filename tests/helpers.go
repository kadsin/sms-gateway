package tests

import (
	"strconv"
	"time"

	"github.com/kadsin/sms-gateway/database/models"
	"github.com/kadsin/sms-gateway/internal/container"
)

func CreateUser(balance ...float32) models.User {
	var b float32 = 100000.55
	if len(balance) > 0 {
		b = balance[0]
	}

	user := models.User{
		Email:   "test" + strconv.Itoa(time.Now().Nanosecond()) + "@example.com",
		Balance: b,
	}

	container.DB().Create(&user)
	return user
}
