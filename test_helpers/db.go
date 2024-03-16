package test_helpers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"runtime"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
)

type Config struct {
	DatabaseHost           string `env:"DATABASE_HOST"`
	DatabaseUser           string `env:"DATABASE_USER"`
	DatabaseName           string `env:"DATABASE_NAME"`
	DatabasePassword       string `env:"DATABASE_PASSWORD"`
	DatabaseMinConnections int    `env:"DB_MIN_CONNECTIONS"`
	DatabaseMaxConnections int    `env:"DB_MAX_CONNECTIONS"`
}

var EnvConfig Config

func LoadEnvForTest() {
	// makes all fields required if default is not defined
	opts := env.Options{RequiredIfNoDef: true}

	EnvConfig = Config{}
	if err := env.Parse(&EnvConfig, opts); err != nil {
		panic(fmt.Sprintf("unable to parse config : %s ", err))
	}
}

// SetupPostgresConnection establishes a connection to the DB.
func SetupPostgresConnection(ctx context.Context, config Config) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", config.DatabaseHost, config.DatabaseUser, config.DatabaseName, config.DatabasePassword)
	err := db.Connect(dsn, 2, 10)
	if err != nil {
		logger.E(ctx, err, "Failed connecting to database", logger.Field("error", err))
		panic(err)
	}

	logger.I(ctx, "Established connection to database")
}

var (
	_, b, _, _            = runtime.Caller(0)
	relPath               = filepath.Join(filepath.Dir(b), "../../db_migrations/db/structure.sql")
	defaultSqlFilePath, _ = filepath.Abs(relPath)
)

func loadSchemaFile(db *sql.DB, schemaFilePath string) {
	content, err := os.ReadFile(schemaFilePath)
	if err != nil {
		panic(fmt.Errorf("%w - error loading schema file %s", err, schemaFilePath))
	}

	_, err = db.Exec(string(content))
	if err != nil {
		panic(fmt.Errorf("%w - error executing schema file %s", err, schemaFilePath))
	}
}

func ClearDataFromPostgres(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	loadSchemaFile(sqlDb, defaultSqlFilePath)
}
