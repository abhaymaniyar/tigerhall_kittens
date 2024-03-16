package db

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations() {
	var db *sql.DB

	db, err := Get().DB()
	if err != nil {
		panic(err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "./internal/db/migrations"); err != nil {
		panic(err)
	}
}
