package model

import (
	"gorm.io/gorm"
	"time"
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
