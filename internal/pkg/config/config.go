package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr    string `yaml:"addr"`
	DB      `yaml:"db"`
	Google  `yaml:"google"`
	Smtp    `yaml:"smtp"`
	Limiter `yaml:"limiter"`
}

type DB struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxIdleTime  string `yaml:"max-idle-time"`
}

type Google struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

type Limiter struct {
	Rps     float64
	Burst   int
	Enabled bool
}

func Load(path string) *Config {
	var conf Config
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		log.Fatal("couldn't read config")
	}
	return &conf
}
