package config

import (
	"time"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
	DB   struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		Name     string `env:"DB_NAME" envDefault:"postgres"`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		User     string `env:"REDIS_USER" envDefault:""`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
	}
	Auth0 struct {
		Domain string `env:"AUTH0_DOMAIN" envDefault:""`
		Aud    string `env:"AUTH0_AUDIENCE" envDefault:""`
	}
	Timeout time.Duration `env:"TIMEOUT" envDefault:"3.5s"`
}
