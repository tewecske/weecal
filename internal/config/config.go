package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Port              string
	SessionCookieName string
	DatabaseName      string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{
		Port:              os.Getenv("LISTEN_ADDR"),
		SessionCookieName: os.Getenv("SESSION_COOKIE_NAME"),
		DatabaseName:      os.Getenv("DATABASE_NAME"),
	}

	return &cfg
}
