package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID      uuid.UUID `gorm:"default: gen_random_uuid();"`
	Balance float32
	Email   string
}
