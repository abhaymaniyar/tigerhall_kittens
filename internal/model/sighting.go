package model

import (
	"time"

	"github.com/google/uuid"
)

type Sighting struct {
	ID               uint `gorm:"primarykey"`
	TigerID          uint
	ReportedByUserID uuid.UUID
	Lat              float64
	Lon              float64
	Timestamp        time.Time
	ImageURL         string
	Tiger            Tiger `gorm:"foreignKey:TigerID" json:"-"`
	ReportedByUser   User  `gorm:"foreignKey:ReportedByUserID" json:"-"`
}
