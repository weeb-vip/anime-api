package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	AppConfig AppConfig
	DBConfig  DBConfig
}

type AppConfig struct {
	APPName string `default:"anime-api"`
	Port    string `default:"3000"`
	Version string `default:"x.x.x"`
}

type DBConfig struct {
	Name     string
	User     string `default:"root"`
	Password string `required:"true" env:"DBPassword" default:"password"`
	Port     uint   `default:"3306"`
}

func LoadConfigOrPanic() Config {
	var config = Config{}
	configor.Load(&config, "config/config.dev.json")

	return config
}
