package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	HTTPAddr    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	user := getenv("DB_USER", "")
	pass := getenv("DB_PASSWORD", "")
	name := getenv("DB_NAME", "")
	if user == "" || pass == "" || name == "" {
		return nil, fmt.Errorf("set DB_USER, DB_PASSWORD, DB_NAME")
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, pass),
		Host:     getenv("DB_HOST", "localhost") + ":" + getenv("DB_PORT", "5432"),
		Path:     name,
		RawQuery: "sslmode=disable",
	}
	return &Config{
		DatabaseUrl: u.String(),
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
	}, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
