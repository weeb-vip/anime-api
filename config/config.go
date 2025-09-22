package config

import "github.com/jinzhu/configor"

type Config struct {
	AppConfig   AppConfig `env:"APPCONFIG"`
	DBConfig    DBConfig
	RedisConfig RedisConfig
}

type AppConfig struct {
	APPName string `default:"anime-api"`
	Port    int    `env:"PORT" default:"3000"`
	Version string `default:"x.x.x" env:"VERSION"`
	Env     string `default:"development" env:"ENV"`
}

type DBConfig struct {
	Host     string `default:"localhost" env:"DBHOST"`
	DataBase string `default:"weeb" env:"DBNAME"`
	User     string `default:"weeb" env:"DBUSERNAME"`
	Password string `required:"true" env:"DBPASSWORD" default:"mysecretpassword"`
	Port     uint   `default:"3306" env:"DBPORT"`
	SSLMode  string `default:"false" env:"DBSSL"`
}

type RedisConfig struct {
	Host     string `default:"localhost" env:"REDIS_HOST"`
	Port     string `default:"6379" env:"REDIS_PORT"`
	Password string `default:"" env:"REDIS_PASSWORD"`
	DB       int    `default:"0" env:"REDIS_DB"`
	Enabled  bool   `default:"false" env:"CACHE_ENABLED"`

	// Cache TTL configurations (in minutes, except LockTTL which is seconds)
	AnimeDataTTLMinutes int `default:"30" env:"CACHE_ANIME_TTL_MINUTES"`
	EpisodeTTLMinutes   int `default:"15" env:"CACHE_EPISODE_TTL_MINUTES"`
	SeasonTTLMinutes    int `default:"60" env:"CACHE_SEASON_TTL_MINUTES"`
	LockTTLSeconds      int `default:"30" env:"CACHE_LOCK_TTL_SECONDS"`
}

func LoadConfigOrPanic() Config {
	var config = Config{}
	configor.Load(&config, "config/config.dev.json")

	return config
}
