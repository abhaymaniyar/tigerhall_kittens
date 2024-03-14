package model

import "time"

type Sighting struct {
	TigerID          uint
	ReportedByUserID uint
	Lat              float64
	Lon              float64
	Timestamp        time.Time
	ImageURL         string
	Tiger            Tiger `gorm:"foreignKey:TigerID"`
	ReportedByUser   User  `gorm:"foreignKey:ReportedByUserID"`
}
