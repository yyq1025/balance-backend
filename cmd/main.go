package main

import (
	"log"

	"github.com/yyq1025/balance-backend/config"
	"github.com/yyq1025/balance-backend/internal/app"

	"github.com/caarlos0/env/v10"
)

func main() {
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
