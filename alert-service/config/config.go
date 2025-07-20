package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	KafkaBroker string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.BindEnv("kafka.broker", "KAFKA_BROKER")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		KafkaBroker: viper.GetString("kafka.broker"),
	}

	return cfg, nil
}
