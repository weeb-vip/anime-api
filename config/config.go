package config

import "github.com/jinzhu/configor"

type Config struct {
	AppConfig AppConfig
	DBConfig  DBConfig
}

type AppConfig struct {
	APPName string `default:"anime-api"`
	Port    int    `env:"PORT" default:"3000"`
	Version string `default:"x.x.x"`
}

type DBConfig struct {
	Host     string `default:"localhost" env:"DBHOST"`
	DataBase string `default:"weeb" env:"DBNAME"`
	User     string `default:"weeb" env:"DBUser"`
	Password string `required:"true" env:"DBPassword" default:"mysecretpassword"`
	Port     uint   `default:"3306" env:"DBPort"`
}

func LoadConfigOrPanic() Config {
	var config = Config{}
	configor.Load(&config, "config/config.dev.json")

	return config
}
