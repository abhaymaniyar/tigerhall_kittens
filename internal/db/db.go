package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"tigerhall_kittens/internal/model"
)

var db *gorm.DB

func Connect(dsn string, maxIdleConnections, maxOpenConnections int) error {
	var err error
	var gormdb *sql.DB

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// Run migrations
	// TODO: use a better db migration tool, ex: goose
	err = db.AutoMigrate(&model.User{}, &model.Tiger{}, &model.Sighting{})
	if err != nil {
		return err
	}

	gormdb, err = db.DB()
	if err != nil {
		return err
	}

	err = gormdb.Ping()
	if err != nil {
		return err
	}

	//gormdb
	gormdb.SetMaxIdleConns(maxIdleConnections)
	gormdb.SetMaxOpenConns(maxOpenConnections)

	return nil
}

func Get() *gorm.DB {
	return db
}

// Close closes the database
func Close() {
	sqldb, _ := db.DB()
	_ = sqldb.Close()
}
