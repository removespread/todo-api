package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	PostgresDSN string `mapstructure:"POSTGRES_DSN"`
	RedisDSN    string `mapstructure:"REDIS_DSN"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
