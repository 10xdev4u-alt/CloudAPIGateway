package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig  `mapstructure:"server"`
	Log    LogConfig     `mapstructure:"log"`
	Routes []RouteConfig `mapstructure:"routes"`
}

type RouteConfig struct {
	Method      string `mapstructure:"method"`
	Path        string `mapstructure:"path"`
	Target      string `mapstructure:"target"`
	StripPrefix bool   `mapstructure:"strip_prefix"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("log.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
