package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect(dsn string, maxIdleConnections, maxOpenConnections int) error {
	var err error
	var gormdb *sql.DB

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
