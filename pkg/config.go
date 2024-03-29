package pkg

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Host         string `mapstructure:"DB_HOST"`
	Port         int    `mapstructure:"DB_PORT"`
	User         string `mapstructure:"DB_USER"`
	Password     string `mapstructure:"DB_PASSWORD"`
	DBname       string `mapstructure:"DB_NAME"`
  RedirectUri  string `mapstructure:"REDIRECT_URI"`
	ClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	ClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	CookieSecret string `mapstructure:"COOKIE_SECRET"`
	State        string `mapstructure:"STATE"`
	ApiNinjas    string `mapstructure:"API_NINJAS"`
}

func Load() (config Config) {
	viper.SetConfigName("development")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	return config
}
