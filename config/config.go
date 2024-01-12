package config

import (
	"time"
)

type Config struct {
	DB struct {
		Dsn      string `env:"DB_DSN" envDefault:""`
	}
	Redis struct {
		Host     string `env:"REDIS_HOST" envDefault:"localhost"`
		Port     string `env:"REDIS_PORT" envDefault:"6379"`
		User     string `env:"REDIS_USER" envDefault:""`
		Password string `env:"REDIS_PASSWORD" envDefault:""`
	}
	Auth0 struct {
		Domain string `env:"AUTH0_DOMAIN,required"`
		Aud    string `env:"AUTH0_AUDIENCE,required"`
	}
	Timeout time.Duration `env:"TIMEOUT" envDefault:"3.5s"`
}
