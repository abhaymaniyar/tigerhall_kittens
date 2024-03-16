package model

import (
	"time"

	"github.com/google/uuid"
)

type Sighting struct {
	ID               uuid.UUID `gorm:"primarykey"`
	TigerID          uint
	ReportedByUserID uuid.UUID
	Lat              float64
	Lon              float64
	SightedAt        time.Time
	ImageURL         string
}
