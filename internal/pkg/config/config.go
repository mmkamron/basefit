package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Addr   string `yaml:"addr"`
	DB     `yaml:"db"`
	Google `yaml:"google"`
	Smtp   `yaml:"smtp"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
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

func Load(path string) *Config {
	var conf Config
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		log.Fatal("couldn't read config")
	}
	return &conf
}
