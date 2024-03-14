package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"tigerhall_kittens/internal/model"
)

var db *gorm.DB

func ConnectAndMigrate() *gorm.DB {
	// TODO: use env variables instead of hard coded strings
	dsn := "host=localhost user=abhay dbname=tigerhall_kittens sslmode=disable password=rootroot"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// Run migrations
	// TODO: use a better db migration tool, ex: goose
	db.AutoMigrate(&model.User{}, &model.Tiger{}, &model.Sighting{})

	return db
}

func Get() *gorm.DB {
	return db
}
