package model

import (
	"time"
)

type Tiger struct {
	ID                uint      `gorm:"primarykey" json:"-"`
	Name              string    `json:"name"  validate:"required"`
	DateOfBirth       time.Time `json:"date_of_birth"  validate:"required"`
	LastSeenTimestamp time.Time `json:"last_seen_at"  validate:"required"`
	LastSeenLat       float64   `json:"last_seen_lat"  validate:"required"`
	LastSeenLon       float64   `json:"last_seen_lon"  validate:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
