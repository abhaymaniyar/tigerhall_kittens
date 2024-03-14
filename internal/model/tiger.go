package model

import (
	"time"
)

type Tiger struct {
	Name              string    `json:"name"  validate:"required"`
	DateOfBirth       time.Time `json:"dateOfBirth"  validate:"required"`
	LastSeenTimestamp time.Time `json:"lastSeenTimestamp"  validate:"required"`
	LastSeenLat       float64   `json:"lastSeenLat"  validate:"required"`
	LastSeenLon       float64   `json:"lastSeenLon"  validate:"required"`
}
