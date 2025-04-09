package main

import (
	"flag"
	"fmt"

	"formy.fprzg.net/internal/controllers"
	"formy.fprzg.net/internal/models"
	"formy.fprzg.net/internal/types"

	_ "github.com/mattn/go-sqlite3"
)

var (
	cfg = types.AppConfig{}
)

func main() {
	flag.StringVar(&cfg.Port, "port", ":3000", "API server port.")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (testing | development | staging | production)")

	flag.StringVar(&cfg.DBDir, "dbDir", "./app.db", "Database directory.")

	// rate-limiter config
	// smtp server config

	flag.Parse()

	if cfg.Env != "development" && cfg.Env != "staging" && cfg.Env != "production" {
		panic(fmt.Errorf("invalid environment: '%s'", cfg.Env))
	}

	m, err := models.GetModels(cfg)
	if err != nil {
		panic(err)
	}

	c := controllers.New(m)

	c.Start(cfg)
}
