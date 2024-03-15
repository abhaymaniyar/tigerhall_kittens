package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strconv"
	"strings"
	"tigerhall_kittens/internal/db"
	"tigerhall_kittens/internal/logger"
)

const (
	EnvDevelopment      = "development"
	EnvTest             = "test"
	EnvStaging          = "staging"
	EnvUnicorn          = "unicorn"
	EnvSandbox          = "sandbox"
	EnvProduction       = "production"
	EnvQualityAssurance = "qa"
	EnvT1Dev            = "t1dev"
	EnvT2Dev            = "t2dev"
)

var Env Config

type Config struct {
	Environment              string `mapstructure:"ENV"`
	Port                     string `mapstructure:"PORT"`
	DatabaseHost             string `mapstructure:"DATABASE_HOST"`
	DatabaseUser             string `mapstructure:"DATABASE_USER"`
	DatabaseName             string `mapstructure:"DATABASE_NAME"`
	DatabasePassword         string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseMinConnections   string `mapstructure:"DB_MIN_CONNECTIONS"`
	DatabaseMaxConnections   string `mapstructure:"DB_MAX_CONNECTIONS"`
	EventPublisherURL        string `mapstructure:"EVENT_PUBLISHER_URL"`
	RedisHost                string `mapstructure:"REDIS_HOST"`
	RedisDB                  string `mapstructure:"REDIS_DB"`
	KafkaUserName            string `mapstructure:"KAFKA_USERNAME"`
	KafkaBrokerList          string `mapstructure:"KAFKA_BROKER_LIST"`
	KafkaConsumerGroupID     string `mapstructure:"KAFKA_CONSUMER_GROUP_ID"`
	KafkaBatchSize           int    `mapstructure:"KAFKA_BATCH_SIZE"`
	KafkaTopicGrowthEvents   string `mapstructure:"KAFKA_TOPIC_GROWTH_EVENTS"`
	KafkaTopicCheckoutEvents string `mapstructure:"KAFKA_TOPIC_CHECKOUT_EVENTS"`
	KafkaTopicShopifyEvents  string `mapstructure:"KAFKA_TOPIC_SHOPIFY_EVENTS"`
	KafkaTopicOmsEvents      string `mapstructure:"KAFKA_TOPIC_OMS_EVENTS"`
	SentryDSN                string `mapstructure:"SENTRY_DSN"`
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
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

	logger.I(ctx, "Established connection to database")
}

func SetupLogger(env string) {
	switch env {
	case EnvDevelopment:
		logger.Init(logger.DEBUG)
	case EnvTest:
		logger.Init(logger.DEBUG)
	case EnvStaging:
		logger.Init(logger.INFO)
		fallthrough
	case EnvUnicorn:
		logger.Init(logger.INFO)
		fallthrough
	case EnvSandbox:
		logger.Init(logger.DEBUG)
		fallthrough
	case EnvProduction:
		logger.Init(logger.INFO)
	}
}
