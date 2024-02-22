package config

import "github.com/jinzhu/configor"

type Config struct {
	AppConfig     AppConfig
	DBConfig      DBConfig
	DataDogConfig DataDogConfig
}

type AppConfig struct {
	APPName string `default:"anime-api"`
	Port    int    `env:"PORT" default:"3000"`
	Version string `default:"x.x.x"`
}

type DBConfig struct {
	Host     string `default:"localhost" env:"DBHOST"`
	DataBase string `default:"weeb" env:"DBNAME"`
	User     string `default:"weeb" env:"DBUSERNAME"`
	Password string `required:"true" env:"DBPASSWORD" default:"mysecretpassword"`
	Port     uint   `default:"3306" env:"DBPORT"`
	SSLMode  string `default:"false" env:"DBSSL"`
}

type DataDogConfig struct {
	DD_AGENT_HOST string `env:"DD_AGENT_HOST" default:"localhost"`
	DD_AGENT_PORT int    `env:"DD_AGENT_PORT" default:"8125"`
}

func LoadConfigOrPanic() Config {
	var config = Config{}
	configor.Load(&config, "config/config.dev.json")

	return config
}
