package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB_SOURCE      string `mapstructure:"DB_SOURCE"`
	BROKER_ADDRESS string `mapstructure:"BROKER_ADDRESS"`
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
