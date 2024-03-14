package model

import (
	"time"
)

type Tiger struct {
	Name              string
	DateOfBirth       time.Time
	LastSeenTimestamp time.Time
	LastSeenLat       float64
	LastSeenLon       float64
}
