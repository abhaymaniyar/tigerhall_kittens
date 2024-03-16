package model

import (
	"time"

	"gorm.io/gorm"
)
import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Username  string    `gorm:"unique"`
	Password  string
	Email     string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
