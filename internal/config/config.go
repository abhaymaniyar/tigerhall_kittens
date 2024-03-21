package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"

	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
)

const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvTest        = "test"
	EnvProduction  = "production"
)

var Env Config

type Config struct {
	Environment            string `mapstructure:"ENV"`
	Port                   string `mapstructure:"PORT"`
	SecretKey              string `mapstructure:"SECRET_KEY"`
	DatabaseHost           string `mapstructure:"DATABASE_HOST"`
	DatabaseUser           string `mapstructure:"DATABASE_USER"`
	DatabaseName           string `mapstructure:"DATABASE_NAME"`
	DatabasePassword       string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseMinConnections string `mapstructure:"DB_MIN_CONNECTIONS"`
	DatabaseMaxConnections string `mapstructure:"DB_MAX_CONNECTIONS"`
}

func bindEnvs(iface interface{}) {
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		err := viper.BindEnv(tv)
		if err != nil {
			logger.E(context.Background(), err, "unable to parse env variables", logger.Field("field", t))
			panic(err)
		}
	}
}

func LoadEnv() error {
	// in case of ENV not set or set to development, read development.env
	if os.Getenv("ENV") == EnvDevelopment || os.Getenv("ENV") == "" {
		viper.SetConfigFile("development.env")

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}

	viper.AutomaticEnv()
	bindEnvs(Env)

	return viper.Unmarshal(&Env)
}

// SetupDBConnection establishes a connection to the DB.
func SetupDBConnection(ctx context.Context) {
	minConnections, err := strconv.Atoi(Env.DatabaseMinConnections)
	if err != nil {
		logger.E(ctx, err, "Invalid DB_MIN_CONNECTIONS", logger.Field("error", err.Error()))
		panic(err)
	}

	maxConnections, err := strconv.Atoi(Env.DatabaseMaxConnections)
	if err != nil {
		logger.E(ctx, err, "Invalid DB_MAX_CONNECTIONS", logger.Field("error", err.Error()))
		panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", Env.DatabaseHost, Env.DatabaseUser, Env.DatabaseName, Env.DatabasePassword)
	err = db.Connect(dsn, minConnections, maxConnections)
	if err != nil {
		logger.E(ctx, err, "Failed connecting to database", logger.Field("error", err))
		panic(err)
	}
	db.RunMigrations()

	logger.I(ctx, "Established connection to database")
}

func SetupLogger(env string) {
	switch env {
	case EnvDevelopment:
		logger.Init(logger.DEBUG)
	case EnvTest, EnvStaging, EnvProduction:
		logger.Init(logger.INFO)
	}
}
