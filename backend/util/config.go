package util

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_SOURCE              string        `mapstructure:"DB_SOURCE"`
	MIGRATION_URL          string        `mapstructure:"MIGRATION_URL"`
	HTTP_SERVER_ADDRESS    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	BROKER_ADDRESS         string        `mapstructure:"BROKER_ADDRESS"`
	REDIS_ADDRESS          string        `mapstructure:"REDIS_ADDRESS"`
	TOKEN_SYMMETRIC_KEY    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	ACCESS_TOKEN_DURATION  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	REFRESH_TOKEN_DURATION time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) *Config {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed reading config file: %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("failed unmarshal config %v", err)
	}

	return &config
}
